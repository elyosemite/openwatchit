package loki

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type lokiResponse struct {
	Status string   `json:"status"`
	Data   lokiData `json:"data"`
}

type lokiData struct {
	ResultType string       `json:"resultType"`
	Result     []lokiStream `json:"result"`
}

type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"` // each element: [nanosecond_timestamp, log_line]
}

type lokiClient struct {
	baseURL    string
	httpClient *http.Client
}

func newLokiClient(baseURL string) *lokiClient {
	return &lokiClient{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *lokiClient) queryRange(ctx context.Context, logql string, start, end time.Time, limit int32) ([]lokiStream, error) {
	params := url.Values{}
	params.Set("query", logql)
	params.Set("start", strconv.FormatInt(start.UnixNano(), 10))
	params.Set("end", strconv.FormatInt(end.UnixNano(), 10))
	if limit > 0 {
		params.Set("limit", strconv.Itoa(int(limit)))
	}

	endpoint := c.baseURL + "/loki/api/v1/query_range?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	applyAuth(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("loki http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("loki returned status %d", resp.StatusCode)
	}

	var lokiResp lokiResponse
	if err := json.NewDecoder(resp.Body).Decode(&lokiResp); err != nil {
		return nil, fmt.Errorf("decode loki response: %w", err)
	}

	return lokiResp.Data.Result, nil
}

// applyAuth reads credentials from environment variables (ADR-0009).
// LOKI_USERNAME + LOKI_PASSWORD → HTTP Basic Auth (Grafana Cloud)
// LOKI_TOKEN                    → Bearer token
func applyAuth(req *http.Request) {
	if username := os.Getenv("LOKI_USERNAME"); username != "" {
		req.SetBasicAuth(username, os.Getenv("LOKI_PASSWORD"))
		return
	}
	if token := os.Getenv("LOKI_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
}
