package loki

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	apiv1 "github.com/elyosemite/openwatchit/api/v1"
)

const version = "0.1.0"

// Server implements apiv1.PluginServiceServer for Grafana Loki.
type Server struct {
	apiv1.UnimplementedPluginServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Capabilities(_ context.Context, _ *apiv1.CapabilitiesRequest) (*apiv1.CapabilitiesResponse, error) {
	return &apiv1.CapabilitiesResponse{
		SignalTypes: []apiv1.SignalType{apiv1.SignalType_SIGNAL_TYPE_LOG},
		Version:     version,
	}, nil
}

func (s *Server) Execute(req *apiv1.ExecuteRequest, stream grpc.ServerStreamingServer[apiv1.NormalizedResultRow]) error {
	baseURL := req.GetConfig().GetParams()["url"]
	if baseURL == "" {
		return fmt.Errorf("missing required backend param: url")
	}

	logql, err := toLogQL(req.GetAst())
	if err != nil {
		return fmt.Errorf("translate AST to LogQL: %w", err)
	}

	now := time.Now()
	start, err := resolveStart(req.GetAst().GetTimeRange().GetLast(), now)
	if err != nil {
		return fmt.Errorf("resolve time range: %w", err)
	}

	streams, err := newLokiClient(baseURL).queryRange(stream.Context(), logql, start, now, req.GetAst().GetLimit())
	if err != nil {
		return fmt.Errorf("loki query: %w", err)
	}

	for _, ls := range streams {
		for _, value := range ls.Values {
			row := normalize(ls, value)
			if row == nil {
				continue
			}
			if err := stream.Send(row); err != nil {
				return err
			}
		}
	}

	return nil
}

// resolveStart converts a relative duration string (e.g. "30m") to an absolute start time.
// Defaults to 1 hour when last is empty.
func resolveStart(last string, now time.Time) (time.Time, error) {
	if last == "" {
		return now.Add(-time.Hour), nil
	}
	d, err := time.ParseDuration(last)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid duration %q: %w", last, err)
	}
	return now.Add(-d), nil
}
