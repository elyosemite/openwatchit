# OpenWatchIt — Domain Glossary

## Signal Type

One of four categories of observability telemetry: **logs**, **metrics**, **traces**, **profiles**. Signal type is declared as the first token in every OWL query (e.g. `logs | ...`, `metrics | ...`). It is not a filter — it defines the shape of the data the query operates on.

## OWL (OpenWatch Query Language)

The unified, pipeline-based query language of OpenWatchIt. KQL-inspired syntax using `|` as the pipe operator. An OWL query targets exactly one Signal Type. The engine translates OWL into each backend's native language; the user never writes PromQL, LogQL, or DogStatsD queries directly.

## Backend

A configured vendor instance. One entry in the TOML config file. A single vendor (e.g. DataDog) may appear as multiple backends (e.g. `datadog-prod`, `datadog-staging`). Each backend declares which Signal Types it supports via its plugin's `Capabilities` response.

## Plugin

An independent process that bridges the OWL engine to a specific vendor's API. Communicates with the engine via gRPC. Can be written in any language. Responsible for: declaring capabilities, translating OWL to the vendor's native query language, executing the query, and streaming back normalized results.

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

## Normalizer

The Control Panel component responsible for validating and enforcing the Normalized Result schema as plugin results arrive. Ensures canonical fields are present and correctly typed before results reach the Result Merger or Renderer.
