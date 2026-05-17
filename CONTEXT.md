# OpenWatchIt — Domain Glossary

## Signal Type

One of four categories of observability telemetry: **logs**, **metrics**, **traces**, **profiles**. Signal type is declared as the first token in every OWL query (e.g. `logs | ...`, `metrics | ...`). It is not a filter — it defines the shape of the data the query operates on.

## OWL (OpenWatch Query Language)

The unified, pipeline-based query language of OpenWatchIt. KQL-inspired syntax using `|` as the pipe operator. An OWL query targets exactly one Signal Type. The engine translates OWL into each backend's native language; the user never writes PromQL, LogQL, or DogStatsD queries directly.

**v0.1 operator scope:** `where` (row filter), `last <duration>` (relative time window), `limit N` (result truncation), `summarize` (aggregation: `count()`, `avg()`, `sum()`). Operators deferred to later versions: `project`, `order by`, `extend`, `join` (cross-signal join targets v0.2).

## Backend

A configured vendor instance. One entry in the TOML config file. A single vendor (e.g. DataDog) may appear as multiple backends (e.g. `datadog-prod`, `datadog-staging`). Each backend declares which Signal Types it supports via its plugin's `Capabilities` response.

## Embedded Mode

The default execution model when running `owit query`, `owit tail`, or `owit repl` without a configured server. The Control Panel runs in-process inside the CLI binary. Plugins run as **lazy daemons** — started on first use, reused across subsequent commands, and terminated after an inactivity timeout (default 30s, configurable in TOML). No running server is required.

## Remote Mode

Activated when `server_url` is set in the TOML config or passed via `--server`. The CLI acts as a thin client, forwarding OWL queries to a remote Control Panel over HTTP/gRPC. Plugins and backend config are managed centrally by the server. The natural deployment model for teams sharing backends and credentials.

## Fan-out Executor

The Control Panel component that dispatches an OWL query (as AST) to the target plugins in parallel via gRPC, and ingests their streaming result rows. Internally contains a **result ingestion pipeline** that normalizes each incoming row against the Normalized Result schema (canonical field validation and type enforcement) before forwarding rows to the Result Merger. Slow backends do not block fast ones.

## Control Panel

The engine core of OpenWatchIt. Contains the OWL Parser, Query Planner, Fan-out Executor (with built-in normalization), Result Merger, Renderer, and API. Runs either embedded inside the CLI binary (Embedded Mode) or as a standalone long-running process started via `owit server` (Remote Mode).

## Query Mode vs. Tail Mode

Two distinct execution modes for the CLI. **Query mode** (`owit query`) buffers all results from all backends, sorts by timestamp, and displays a complete ordered result set. **Tail mode** (`owit tail`) streams results continuously as they arrive from backends, in real time, without ordering guarantees — analogous to `tail -f`.

## Plugin

An independent process that bridges the OWL engine to a specific vendor's API. Communicates with the engine via gRPC. Can be written in any language. Responsible for: declaring capabilities, translating an OWL AST to the vendor's native query language, executing the query, and streaming back Normalized Results.

Plugins are installed via `owit plugin install <name>` (to `~/.owit/plugins/`) or `owit plugin install --local <name>` (to `.owit/plugins/` relative to the active TOML config). The CLI resolves local plugins before global ones.

## Plugin Marketplace

The public registry where plugin authors publish versioned vendor integrations. `owit plugin install datadog` resolves the plugin name against this registry and downloads the appropriate platform binary. Private plugins can be installed by name with an org scope (`@mycompany/internal-splunk`) or directly by binary path.

## Fan-out

The act of dispatching a single OWL query to multiple backends in parallel. Backends respond independently; slow backends do not block fast ones.

Default fan-out scope (no explicit flags): all backends where `default = true` **and** whose plugin Capabilities declare support for the query's signal type. Explicit flags (`--backends`, `--tag`, `--type`) override this entirely.

## Query Planner

The Control Panel component that determines the execution plan for an OWL query before dispatch. Responsibilities: resolve fan-out scope (which backends receive the query), detect cross-signal joins and choose between plugin-native or core-side execution (see ADR-0003), and validate that the query's signal type is supported by the targeted backends.

## Normalized Result

The common row schema returned by all plugins regardless of vendor. Contains two layers:

1. **Canonical fields** — a fixed set of semantically named fields every plugin must map to, regardless of what the vendor calls them internally. These are the fields OWL operators (`join`, `where`, `project`) can reference by name across vendors. Examples: `timestamp`, `_source`, `trace_id`, `span_id`, `service`, `level`, `message`.

2. **Passthrough fields** — vendor-specific fields the plugin returns with their native names (e.g. `dd.trace_id`, `xray_trace_id`). Available to the user but carry no cross-vendor uniformity guarantee.

OWL joins and correlations operate only on canonical fields. A future **Plugin Field Mapping** mechanism (not in v0.x) will let plugin authors declare aliases that promote vendor-specific fields into canonical names, enabling users to write queries in a vocabulary closer to their own team's language.

