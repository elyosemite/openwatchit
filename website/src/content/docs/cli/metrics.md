---
title: owit metrics
description: Query metric data across configured backends.
---

Query metric time series across all configured backends that support the `metrics` signal type.

## Usage

```bash
owit metrics [filters] [flags]
```

## Examples

```bash
# All HTTP request metrics
owit metrics name eq http_requests_total --last 1h

# Metrics from a specific environment
owit metrics name eq cpu_usage env=prod --last 30m

# Aggregate: sum by service
owit metrics name eq http_requests_total --summarize "sum(value) by service" --last 1h

# High error rate
owit metrics name eq error_rate value gt 0.05 --last 15m --tag prod

# Specific backends
owit metrics name eq memory_usage --backends prometheus --last 2h
```

## Canonical fields

| Field | Description |
|---|---|
| `timestamp` | RFC 3339 timestamp |
| `_source` | Backend name |
| `service` | Service name |
| `name` | Metric name (maps to `__name__` in Prometheus) |
| `value` | Metric value |

## See also

- [Filter Operators](/owl/operators/)
- [Flags Reference](/owl/flags/)
