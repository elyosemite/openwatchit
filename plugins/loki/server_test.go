package loki_test

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apiv1 "github.com/elyosemite/openwatchit/api/v1"
	lokiplugin "github.com/elyosemite/openwatchit/plugins/loki"
)

func newTestServer(t *testing.T) (apiv1.PluginServiceClient, func()) {
	t.Helper()

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	grpcSrv := grpc.NewServer()
	apiv1.RegisterPluginServiceServer(grpcSrv, lokiplugin.NewServer())
	go grpcSrv.Serve(lis) //nolint:errcheck

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		conn.Close()
		grpcSrv.Stop()
	}

	return apiv1.NewPluginServiceClient(conn), cleanup
}

func lokiMockHandler(streams []map[string]any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"status": "success",
			"data": map[string]any{
				"resultType": "streams",
				"result":     streams,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp) //nolint:errcheck
	}
}

func TestCapabilities(t *testing.T) {
	client, cleanup := newTestServer(t)
	defer cleanup()

	resp, err := client.Capabilities(context.Background(), &apiv1.CapabilitiesRequest{})
	if err != nil {
		t.Fatal(err)
	}

	if len(resp.SignalTypes) != 1 || resp.SignalTypes[0] != apiv1.SignalType_SIGNAL_TYPE_LOG {
		t.Errorf("expected [LOG], got %v", resp.SignalTypes)
	}
	if resp.Version == "" {
		t.Error("expected non-empty version")
	}
}

func TestExecute_StreamsNormalizedRows(t *testing.T) {
	mockLoki := httptest.NewServer(lokiMockHandler([]map[string]any{
		{
			"stream": map[string]string{"level": "error", "service": "checkout"},
			"values": [][]string{
				{"1715097600000000000", "checkout payment failed"},
			},
		},
	}))
	defer mockLoki.Close()

	client, cleanup := newTestServer(t)
	defer cleanup()

	stream, err := client.Execute(context.Background(), &apiv1.ExecuteRequest{
		Ast: &apiv1.AST{
			SignalType: apiv1.SignalType_SIGNAL_TYPE_LOG,
			Filters: []*apiv1.Filter{
				{Field: "level", Op: apiv1.FilterOp_FILTER_OP_EQ, Value: "error"},
			},
			TimeRange: &apiv1.TimeRange{Last: "30m"},
			Limit:     10,
		},
		Config: &apiv1.BackendConfig{
			Name:   "loki-test",
			Params: map[string]string{"url": mockLoki.URL},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	var rows []*apiv1.NormalizedResultRow
	for {
		row, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		rows = append(rows, row)
	}

	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}

	row := rows[0]
	if row.Level != "error" {
		t.Errorf("level: got %q, want %q", row.Level, "error")
	}
	if row.Service != "checkout" {
		t.Errorf("service: got %q, want %q", row.Service, "checkout")
	}
	if row.Message != "checkout payment failed" {
		t.Errorf("message: got %q, want %q", row.Message, "checkout payment failed")
	}
	if row.Source != "loki" {
		t.Errorf("source: got %q, want %q", row.Source, "loki")
	}
	if row.TimestampNs != 1715097600000000000 {
		t.Errorf("timestamp_ns: got %d, want 1715097600000000000", row.TimestampNs)
	}
}

func TestExecute_NeFilter(t *testing.T) {
	var capturedQuery string
	mockLoki := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.Query().Get("query")
		resp := map[string]any{
			"status": "success",
			"data":   map[string]any{"resultType": "streams", "result": []any{}},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp) //nolint:errcheck
	}))
	defer mockLoki.Close()

	client, cleanup := newTestServer(t)
	defer cleanup()

	stream, err := client.Execute(context.Background(), &apiv1.ExecuteRequest{
		Ast: &apiv1.AST{
			SignalType: apiv1.SignalType_SIGNAL_TYPE_LOG,
			Filters: []*apiv1.Filter{
				{Field: "level", Op: apiv1.FilterOp_FILTER_OP_NE, Value: "debug"},
			},
			TimeRange: &apiv1.TimeRange{Last: "1h"},
		},
		Config: &apiv1.BackendConfig{
			Params: map[string]string{"url": mockLoki.URL},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
	}

	want := `{level!="debug"}`
	if capturedQuery != want {
		t.Errorf("LogQL sent to Loki: got %q, want %q", capturedQuery, want)
	}
}

func TestExecute_MissingURL(t *testing.T) {
	client, cleanup := newTestServer(t)
	defer cleanup()

	stream, err := client.Execute(context.Background(), &apiv1.ExecuteRequest{
		Ast:    &apiv1.AST{SignalType: apiv1.SignalType_SIGNAL_TYPE_LOG},
		Config: &apiv1.BackendConfig{Params: map[string]string{}},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = stream.Recv()
	if err == nil {
		t.Fatal("expected error when url is missing, got nil")
	}
}
