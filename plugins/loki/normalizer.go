package loki

import (
	"strconv"

	apiv1 "github.com/elyosemite/openwatchit/api/v1"
)

// canonicalFields maps stream label names that promote to canonical NormalizedResultRow fields.
var canonicalFields = map[string]bool{
	"level":    true,
	"service":  true,
	"trace_id": true,
	"span_id":  true,
}

// normalize converts a single Loki stream + value pair into a NormalizedResultRow.
// Canonical labels (level, service, trace_id, span_id) are promoted to named fields.
// Remaining labels go into Passthrough.
func normalize(stream lokiStream, value []string) *apiv1.NormalizedResultRow {
	if len(value) < 2 {
		return nil
	}

	tsNs, _ := strconv.ParseInt(value[0], 10, 64)

	row := &apiv1.NormalizedResultRow{
		TimestampNs: tsNs,
		Source:      "loki",
		Message:     value[1],
		Level:       stream.Stream["level"],
		Service:     stream.Stream["service"],
		TraceId:     stream.Stream["trace_id"],
		SpanId:      stream.Stream["span_id"],
		Passthrough: make(map[string]string),
	}

	for k, v := range stream.Stream {
		if !canonicalFields[k] {
			row.Passthrough[k] = v
		}
	}

	return row
}
