# ADR-0007: Fan-out returns partial results with explicit per-backend warnings on failure

## Status
Accepted

## Context
During a fan-out query, one or more backends may fail (timeout, auth error, network unreachable). A policy is needed for what the user receives in that case.

## Decision
The engine returns results from all backends that succeeded, and emits a clearly visible warning for each backend that failed, including the backend name and the reason for failure. Warnings appear before or alongside the result set — never silently swallowed.

Example output:
```
⚠ datadog: timeout after 30s
⚠ cloudwatch: authentication failed
──────────────────────────────────────────
[results from loki]
```

A per-backend timeout is configurable in the TOML config. If all backends fail, the output is all warnings and no result rows.

## Reasons
- Silent partial results are dangerous in incident response: an engineer may conclude "no errors found" when backends were simply unreachable.
- Explicit warnings preserve the utility of partial data (the engineer still gets Loki's results) while making the incompleteness impossible to miss.
- Fail-total (Option A) discards useful data unnecessarily — if Loki is healthy, its results have value even if DataDog is down.

## Consequences
- The Fan-out Executor must track per-backend success/failure independently and pass failure metadata to the Renderer alongside result rows.
- The Renderer must emit warnings before the result table, not buried after it.
- JSON/machine-readable output (`--output json`) must include a top-level `warnings` array alongside the `rows` array so scripts can detect partial results programmatically.
