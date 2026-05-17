# ADR-0004: Fan-out scope is determined by default flag + signal type match

## Status
Accepted

## Context
When a user runs `owit query "logs | ..."` without explicit targeting flags, the Query Planner must decide which backends receive the query.

## Decision
The Query Planner fans out to all backends that satisfy both conditions:

1. `default = true` in the TOML configuration file.
2. The backend's plugin declares support for the query's signal type in its `Capabilities` response.

If multiple backends meet both conditions (e.g. Loki and DataDog both default and both supporting logs), the query is dispatched to all of them in parallel.

Explicit flags override this logic entirely:
- `--backends loki,datadog` — send only to named backends, ignoring default flag
- `--tag prod` — send to all backends with that tag that match the signal type
- `--type logs` — send to all backends supporting logs, ignoring default flag

## Reasons
- Simple mental model for the user: "default means it always runs unless I say otherwise."
- Signal type filtering via Capabilities avoids routing errors silently (a metrics-only backend never receives a logs query).
- The TOML config is the single source of truth for what's default; no hidden runtime state.

## Consequences
- The Query Planner must query plugin Capabilities at startup (or cache them) to know which signal types each backend supports.
- Users configuring backends must mark each one as `default = true` or `default = false` explicitly — there is no implicit defaulting.
- A backend that supports multiple signal types (e.g. DataDog: logs + metrics + traces) and is marked default will receive any query whose signal type it supports.
