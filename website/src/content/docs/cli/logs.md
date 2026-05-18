---
title: owit logs
description: Query log data across configured backends.
---

Query log entries across all configured backends that support the `logs` signal type.

## Usage

```bash
owit logs [filters] [flags]
```

## Examples

```bash
# All error logs in the last hour
owit logs level=error --last 1h --limit 100

# Errors from a specific service
owit logs level=error service=payments --last 30m

# Logs containing a specific message
owit logs message contains "connection refused" --last 2h

# Not info level, from production
owit logs level ne info --tag prod --last 1h

# Output as JSON for scripting
owit logs level=error --last 1h --output json | jq '.rows[] | .message'

# Specific backends only
owit logs level=error --backends loki,datadog --last 30m
```

## Canonical fields

| Field | Description |
|---|---|
| `timestamp` | RFC 3339 timestamp |
| `_source` | Backend name |
| `service` | Service name |
| `level` | Log level (`error`, `warn`, `info`, `debug`) |
| `message` | Log line body |
| `trace_id` | Distributed trace identifier (when present) |
| `span_id` | Span identifier (when present) |

## See also

- [Filter Operators](/owl/operators/)
- [Flags Reference](/owl/flags/)
- [owit tail](/cli/tail/) — stream logs in real time
