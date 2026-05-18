---
title: Filter Operators
description: OWL filter operator reference.
---

OWL filters are written as positional triplets directly after the signal type subcommand:

```
owit <signal-type> <field> <op> <value> [<field> <op> <value> ...] [flags]
```

## Equality shorthand

`field=value` is syntactic sugar for `field eq value`:

```bash
owit logs level=error service=checkout --last 30m
# equivalent to:
owit logs level eq error service eq checkout --last 30m
```

## Operator reference

| Operator | Meaning | Example |
|---|---|---|
| `eq` | Equal | `level eq error` |
| `ne` | Not equal | `level ne info` |
| `gt` | Greater than | `duration gt 500ms` |
| `ge` | Greater than or equal | `duration ge 1s` |
| `lt` | Less than | `duration lt 100ms` |
| `le` | Less than or equal | `status_code le 499` |
| `contains` | Substring match | `message contains timeout` |

## Multiple filters

Multiple filter triplets are combined with AND logic:

```bash
owit logs level=error service=payments message contains timeout --last 1h
# Returns rows where level=error AND service=payments AND message contains "timeout"
```

## Combining with flags

Filters always come before flags. Flags start with `--`:

```bash
owit traces duration gt 500ms root_error eq true --last 1h --limit 20 --backends tempo,datadog
```
