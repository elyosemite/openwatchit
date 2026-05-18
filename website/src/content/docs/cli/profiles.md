---
title: owit profiles
description: Query continuous profiling data across configured backends.
---

Query CPU, memory, and other profiling data across all configured backends that support the `profiles` signal type.

## Usage

```bash
owit profiles [filters] [flags]
```

## Examples

```bash
# CPU profiles for a service
owit profiles type=cpu service=api-gateway --last 1h

# Top memory consumers
owit profiles type=memory --summarize "avg(value) by function" --last 30m --limit 10

# High CPU functions
owit profiles type=cpu value gt 80 --last 15m
```

## Canonical fields

| Field | Description |
|---|---|
| `timestamp` | RFC 3339 timestamp |
| `_source` | Backend name |
| `service` | Service name |
| `type` | Profile type (`cpu`, `memory`, `goroutine`, etc.) |
| `function` | Function name |
| `value` | Profiling value (unit depends on type) |

## See also

- [Filter Operators](/owl/operators/)
- [Flags Reference](/owl/flags/)
