# OpenWatchIt — Domain Glossary

## Signal Type

One of four categories of observability telemetry: **logs**, **metrics**, **traces**, **profiles**. Every OWL query targets exactly one signal type, expressed as the CLI subcommand (`owit traces ...`, `owit logs ...`). It is not a filter — it defines the shape of the data the query operates on.

## OWL (OpenWatch Query Language)

The unified query language of OpenWatchIt. OWL is expressed exclusively through CLI flags: the signal type is the subcommand and each operator is a named flag.

```bash
owit traces --where "duration > 1s" --last 1h --limit 20
owit logs --where "level == 'error'" --where "service == 'payments'" --last 30m
```

The engine translates OWL into each backend's native language; the user never writes PromQL, LogQL, or DogStatsD queries directly.

**v0.1 operator scope:** `--where` (row filter), `--last <duration>` (relative time window), `--limit N` (result truncation), `--summarize` (aggregation: `count()`, `avg()`, `sum()`). Operators deferred to later versions: `--project`, `--order-by`, `--extend`, `--join` (cross-signal join targets v0.2).

## Backend

A configured vendor instance. One entry in the TOML config file. A single vendor (e.g. DataDog) may appear as multiple backends (e.g. `datadog-prod`, `datadog-staging`). Each backend declares which Signal Types it supports via its plugin's `Capabilities` response.

## Embedded Mode

The default execution model when running `owit <signal-type>` or `owit tail` without a configured server. The Control Panel runs in-process inside the CLI binary. Plugins run as **lazy daemons** — started on first use, reused across subsequent commands, and terminated after an inactivity timeout (default 30s, configurable in TOML). No running server is required.

## Remote Mode

Activated when `server_url` is set in the TOML config or passed via `--server`. The CLI acts as a thin client, forwarding OWL queries to a remote Control Panel over HTTP/gRPC. Plugins and backend config are managed centrally by the server. The natural deployment model for teams sharing backends and credentials.

## Plugin Manager

The Control Panel component responsible for plugin process lifecycle. Responsibilities: discover installed plugins from `~/.owit/plugins/` and `.owit/plugins/`; spawn plugin processes on demand; maintain gRPC connections; health-check running plugins; enforce the inactivity timeout (ADR-0006); and cache each plugin's `Capabilities` response for use by the Query Planner. The Dispatcher and Query Planner never spawn or connect to plugins directly — they always go through the Plugin Manager.

## Dispatcher

The Control Panel component that dispatches an OWL query (as AST) to the target plugins in parallel via gRPC and ingests their streaming result rows. Internally contains a **result ingestion pipeline** that normalizes each incoming row against the Normalized Result schema (canonical field validation and type enforcement) before forwarding rows to the Result Merger. Slow backends do not block fast ones.

## Control Panel

The engine core of OpenWatchIt. Contains: OWL Parser, Query Planner, Plugin Manager, Dispatcher (with built-in normalization), Result Merger, and API. Runs either embedded inside the CLI binary (Embedded Mode) or as a standalone long-running process started via `owit server` (Remote Mode).

The Control Panel is **not responsible for rendering output**. It only produces structured data (Normalized Result rows + warnings). Rendering is a CLI concern.

## Renderer

A CLI-layer component, not part of the Control Panel. Receives structured result rows and warnings from the engine (directly in Embedded Mode, or via API response in Remote Mode) and formats them for the terminal: table, JSON, or streaming. The Browser UI (React) has its own rendering and never uses this component.

## Query Mode vs. Tail Mode

Two distinct execution modes for the CLI. **Query mode** (`owit <signal-type> [flags]`) buffers all results from all backends, sorts by timestamp, and displays a complete ordered result set. **Tail mode** (`owit tail <signal-type> [flags]`) streams results continuously as they arrive from backends, in real time, without ordering guarantees — analogous to `tail -f`.

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

