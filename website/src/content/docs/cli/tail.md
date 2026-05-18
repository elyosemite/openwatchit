---
title: owit tail
description: Stream observability data in real time.
---

Stream results continuously from configured backends as they arrive — analogous to `tail -f` for logs.

## Usage

```bash
owit tail <signal-type> [filters] [flags]
```

## Examples

```bash
# Stream all incoming error logs
owit tail logs level=error

# Stream errors from a specific service
owit tail logs level=error service=payments

# Stream slow traces in real time
owit tail traces duration gt 500ms

# From specific backends only
owit tail logs level=error --backends loki,datadog
```

## Behaviour

- Results appear as they arrive from each backend — no waiting for all backends to respond.
- Results are **not** sorted by timestamp. Rows from fast backends appear before rows from slow backends.
- The stream continues until interrupted with `Ctrl+C`.
- Backend failures appear as inline warnings and do not stop the stream.

## Difference from `owit logs`

| | `owit logs` | `owit tail` |
|---|---|---|
| Result order | Sorted by timestamp | Arrival order |
| Completes | Yes (buffered) | No (continuous stream) |
| Use case | Querying past data | Following live activity |

## See also

- [owit logs](/cli/logs/)
- [Filter Operators](/owl/operators/)
- [Flags Reference](/owl/flags/)
