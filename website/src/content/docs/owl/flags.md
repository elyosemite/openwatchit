---
title: Flags Reference
description: Named flags available on all OWL commands.
---

Named flags are available on all signal type subcommands (`owit logs`, `owit metrics`, `owit traces`, `owit profiles`) and on `owit tail`.

## `--last <duration>`

Relative time window. Accepted units: `s` (seconds), `m` (minutes), `h` (hours), `d` (days).

```bash
owit logs level=error --last 30m
owit metrics name eq http_requests_total --last 2h
owit traces duration gt 1s --last 7d
```

## `--limit <n>`

Maximum number of result rows returned.

```bash
owit logs level=error --last 1h --limit 100
```

## `--summarize "<expr>"`

Aggregation expression. Supported functions: `count()`, `avg(<field>)`, `sum(<field>)`.

```bash
owit metrics name eq http_requests_total --summarize "sum(value) by service" --last 1h
owit logs level=error --summarize "count() by service" --last 24h
```

## `--backends <name>[,<name>...]`

Target specific backends by name, bypassing the default fan-out scope.

```bash
owit logs level=error --backends loki,datadog --last 1h
```

## `--tag <tag>`

Target all backends with a given tag, bypassing the default fan-out scope.

```bash
owit logs level=error --tag prod --last 1h
owit metrics name eq cpu_usage --tag staging --last 30m
```

## `--output <format>`

Output format. Options: `table` (default), `json`.

```bash
owit logs level=error --last 1h --output json | jq '.rows[] | .message'
```

In `json` mode the response includes both `rows` and `warnings` arrays, enabling scripts to detect partial results programmatically.
