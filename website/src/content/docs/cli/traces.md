---
title: owit traces
description: Query trace data across configured backends.
---

Query distributed traces across all configured backends that support the `traces` signal type.

## Usage

```bash
owit traces [filters] [flags]
```

## Examples

```bash
# Slow traces
owit traces duration gt 500ms --last 1h --limit 20

# Root errors only
owit traces root_error eq true --last 30m

# Slow errors from a specific service
owit traces duration gt 1s root_error eq true service=api-gateway --last 1h

# Specific trace ID
owit traces trace_id eq abc123def456 --last 24h
```

## Canonical fields

| Field | Description |
|---|---|
| `timestamp` | RFC 3339 start time of the trace |
| `_source` | Backend name |
| `service` | Root service name |
| `trace_id` | Distributed trace identifier |
| `span_id` | Root span identifier |
| `duration` | Total trace duration (e.g. `1.2s`, `450ms`) |
| `root_error` | `true` if the root span ended in error |

## See also

- [Filter Operators](/owl/operators/)
- [Flags Reference](/owl/flags/)
