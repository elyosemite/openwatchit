---
title: gRPC Contract
description: Full reference for the ObservabilityPlugin gRPC service.
---

The OpenWatchIt plugin contract is defined in three protobuf files in the repository under `proto/`.

## Files

| File | Contents |
|---|---|
| `proto/ast.proto` | `QueryAST`, `FilterNode`, `SignalType`, `FilterOperator` |
| `proto/result.proto` | `ResultRow`, `QueryWarning` |
| `proto/plugin.proto` | `ObservabilityPlugin` service and all request/response messages |

## Service definition

```proto
service ObservabilityPlugin {
  rpc Capabilities(CapabilitiesRequest) returns (CapabilitiesResponse);
  rpc Translate(TranslateRequest)       returns (TranslateResponse);
  rpc Query(QueryRequest)               returns (stream QueryResult);
  rpc HealthCheck(HealthRequest)        returns (HealthResponse);
}
```

## QueryAST

The structured representation of an OWL query passed to plugins via `Translate` and `Query`. Plugins receive this — never raw OWL text.

| Field | Type | Description |
|---|---|---|
| `signal_type` | `SignalType` | The signal type being queried |
| `filters` | `repeated FilterNode` | All filter triplets |
| `time_window` | `TimeWindow` | Resolved `--last` duration in seconds |
| `limit` | `uint32` | Value of `--limit` (0 = no limit) |
| `summarize` | `SummarizeExpr` | Value of `--summarize` |

## FilterNode

| Field | Type | Description |
|---|---|---|
| `field` | `string` | Field name (e.g. `level`, `service`, `duration`) |
| `operator` | `FilterOperator` | One of `EQ`, `NE`, `GT`, `GE`, `LT`, `LE`, `CONTAINS` |
| `value` | `string` | Filter value as string |

## ResultRow

| Field | Type | Description |
|---|---|---|
| `timestamp` | `string` | RFC 3339 timestamp |
| `source` | `string` | Backend name from TOML config |
| `service` | `string` | Service name (canonical) |
| `level` | `string` | Log level (canonical, logs only) |
| `message` | `string` | Log message body (canonical, logs only) |
| `trace_id` | `string` | Trace identifier (canonical, when present) |
| `span_id` | `string` | Span identifier (canonical, when present) |
| `passthrough` | `map<string, string>` | Vendor-specific fields, native names |

## Versioning

Breaking changes to the protobuf schema require a major version bump. The package is versioned as `openwatchit.v1`. Plugin authors should pin to a major version.
