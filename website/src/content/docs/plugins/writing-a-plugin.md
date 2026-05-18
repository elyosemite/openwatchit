---
title: Writing a Plugin
description: How to build a vendor integration plugin for OpenWatchIt.
---

A plugin is an independent process that bridges the OpenWatchIt engine to a specific vendor's API. Plugins communicate with the engine via gRPC and can be written in any language that supports gRPC — Go, Rust, Python, TypeScript, Java, and more.

## How it works

1. The engine spawns your plugin binary as a child process.
2. Your plugin starts a gRPC server on a local port.
3. The engine calls `Capabilities` once and caches the response.
4. For each query, the engine calls `Translate` (AST → vendor query string) then `Query` (execute and stream results).
5. Your plugin streams back `ResultRow` messages in the Normalized Result schema.

## The plugin contract

Your plugin must implement the `ObservabilityPlugin` gRPC service defined in `proto/plugin.proto`:

```proto
service ObservabilityPlugin {
  rpc Capabilities(CapabilitiesRequest) returns (CapabilitiesResponse);
  rpc Translate(TranslateRequest)       returns (TranslateResponse);
  rpc Query(QueryRequest)               returns (stream QueryResult);
  rpc HealthCheck(HealthRequest)        returns (HealthResponse);
}
```

See [gRPC Contract](/plugins/grpc-contract/) for the full schema.

## Capabilities

Declare what your plugin supports. The Query Planner uses this to route queries correctly:

```go
func (p *Plugin) Capabilities(_ context.Context, _ *pb.CapabilitiesRequest) (*pb.CapabilitiesResponse, error) {
    return &pb.CapabilitiesResponse{
        PluginName:      "loki",
        PluginVersion:   "0.1.0",
        SignalTypes:     []pb.SignalType{pb.SignalType_SIGNAL_TYPE_LOGS},
        FilterOperators: []pb.FilterOperator{
            pb.FilterOperator_FILTER_OPERATOR_EQ,
            pb.FilterOperator_FILTER_OPERATOR_NE,
            pb.FilterOperator_FILTER_OPERATOR_CONTAINS,
        },
    }, nil
}
```

## Translate

Receive the OWL AST and return the vendor-native query string. Your plugin never parses OWL text — only walks the AST nodes:

```go
func (p *Plugin) Translate(_ context.Context, req *pb.TranslateRequest) (*pb.TranslateResponse, error) {
    ast := req.Ast
    // Build a LogQL query from ast.Filters, ast.TimeWindow, etc.
    query := buildLogQL(ast)
    return &pb.TranslateResponse{NativeQuery: query}, nil
}
```

## Emulating unsupported operators

If your vendor does not natively support an operator (e.g. Prometheus has no `LIMIT`), fetch the full result set and apply the operation in memory before streaming rows back. Declare only natively-translated operators in `Capabilities.FilterOperators`.

## Normalized Result

Every `ResultRow` must populate the canonical fields. Vendor-specific fields go in `passthrough`:

```go
row := &pb.ResultRow{
    Timestamp: entry.Timestamp.Format(time.RFC3339),
    Source:    "production-loki",
    Service:   entry.Labels["app"],
    Level:     entry.Labels["level"],
    Message:   entry.Line,
    Passthrough: map[string]string{
        "loki.stream": entry.Labels["stream"],
    },
}
```

## Installing your plugin

Place the compiled binary in `~/.owit/plugins/<name>` and reference it in your TOML config:

```toml
[backends.my-loki]
plugin = "loki"
url    = "https://loki.internal"
```

## Publishing to the Plugin Marketplace

_(Coming in a future release.)_
