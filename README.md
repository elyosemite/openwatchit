# OpenWatchIt — One Query Language to Rule Them All

[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Built with Go](https://img.shields.io/badge/built%20with-Go-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![gRPC](https://img.shields.io/badge/plugins-gRPC-244c5a?logo=grpc&logoColor=white)](https://grpc.io)
[![GitHub Discussions](https://img.shields.io/github/discussions/elyosemite/openwatchit?label=discussions&color=6e40c9&logo=github)](https://github.com/elyosemite/openwatchit/discussions)
[![GitHub Stars](https://img.shields.io/github/stars/elyosemite/openwatchit?style=flat&logo=github&color=f5a623)](https://github.com/elyosemite/openwatchit/stargazers)
[![Contributions Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg)](https://github.com/elyosemite/openwatchit/discussions)

---

## The Problem

Modern software systems are observable — but at a cost nobody talks about enough.

Your microservices run on AWS. Your traces live in Datadog. Your logs are split between Loki and CloudWatch. Your team recently acquired a company that runs everything on Azure with Application Insights. Your on-call engineer wakes up at 3am and has to open four browser tabs, remember four query syntaxes, and mentally join the results across four completely different UIs — all while a production incident is burning.

This is the state of observability in 2026 for most engineering teams. Not because people chose it, but because it grew that way.

The tools are excellent individually. The problem is **the seams between them**.

---

## What OpenWatchIt Is

**OpenWatchIt** is an open-source observability gateway — a CLI, a UI, and a plugin runtime — that lets you query logs, metrics, traces, and profiles across any combination of vendors using a **single, unified query language**.

You write one query. OpenWatchIt dispatches it to every configured backend in parallel, translates it to each vendor's native language under the hood, and returns a unified, normalized result.

```bash
owit logs --where "level == 'error'" --where "service == 'payments'" --last 1h --limit 50
```

That single command hits Loki, Datadog, and CloudWatch simultaneously — you get one result, ranked by timestamp, with a `_source` column showing where each entry came from.

No new vendor to learn. No new dashboard to configure. One tool, every source.

---

## Core Pillars

### 1. OpenWatch Query Language (OWL)

OWL is OpenWatchIt's own query language. The signal type is the subcommand; each operator is a named flag:

```bash
# Logs
owit logs --where "level == 'error'" --where "service == 'checkout'" --last 30m --limit 100

# Metrics
owit metrics --where "__name__ == 'http_requests_total'" --where "env == 'prod'" --summarize "sum(value) by service"

# Traces
owit traces --where "duration > 500ms" --where "root_error == true" --limit 20

# Profiles
owit profiles --where "type == 'cpu'" --where "service == 'api-gateway'" --summarize "avg(value) by function" --limit 10
```

OWL is intentionally readable. A developer who has never used it before should be able to write a useful query in under five minutes.

### 2. Multi-vendor Dispatch

Queries are dispatched to multiple backends in parallel. By default, a query is sent to every backend marked `default = true` in your config **that declares support for the query's signal type** (logs, metrics, traces, or profiles). A Prometheus backend never receives a logs query; a Loki backend never receives a metrics query.

```bash
# All default backends matching the signal type
owit logs --where "level == 'error'"

# Specific backends
owit logs --where "service == 'api'" --backends loki,datadog

# All backends of a given type, regardless of default flag
owit logs --where "level == 'error'" --type logs

# All backends with a given tag
owit metrics --where "__name__ == 'http_requests_total'" --tag prod
```

Results arrive as backends respond. Slow backends don't block fast ones. If a backend fails, you get results from the healthy backends plus an explicit warning — never silent partial data.

### 3. Plugin Architecture via gRPC

Every vendor integration is a **plugin** — an independent process that speaks a gRPC contract defined by OpenWatchIt core. Plugins can be written in any language: Go, Rust, Python, TypeScript, Java — anything that can implement a gRPC server.

```proto
service ObservabilityPlugin {
  rpc Capabilities(CapabilitiesRequest) returns (CapabilitiesResponse);
  rpc Translate(TranslateRequest) returns (TranslateResponse);
  rpc Query(QueryRequest) returns (stream QueryResult);
  rpc Validate(ValidateRequest) returns (ValidateResponse);
  rpc HealthCheck(HealthRequest) returns (HealthResponse);
}
```

The core OWL Parser produces an **AST (Abstract Syntax Tree)** from the user's query. `TranslateRequest` carries this AST — not raw OWL text. Plugins never implement an OWL parser; they only walk the AST and emit their vendor's native query language (PromQL, LogQL, DogStatsD, KQL, etc.). This means OWL syntax can evolve without breaking existing plugins.

This means:
- The core team maintains the language and the runtime.
- The community maintains vendor integrations.
- Enterprise teams can write private plugins for internal tools without forking anything.
- Plugins are versioned, sandboxed, and loaded dynamically.

### 4. Plugin Marketplace

A community registry where plugin authors publish, version, and document their integrations:

```bash
owit plugin install datadog
owit plugin install @mycompany/internal-splunk
owit plugin install tempo --version 2.1.0
owit plugin list
owit plugin update --all
```

Any developer can publish a plugin. Plugins go through a lightweight verification process (schema validation, basic security scan) before appearing in the public registry.

### 5. CLI-first, UI when you need it

**CLI** — designed for engineers who live in the terminal, scripts, CI pipelines, and incident response playbooks:

```bash
owit traces --where "duration > 1s" --output table
owit metrics --where "env == 'prod'" --output json | jq '.[] | .value'
owit tail logs --where "service == 'api'"          # streaming, like tail -f
```

**UI** — a browser-based interface that connects to an OpenWatchIt server instance, designed for teams who want a shared, visual layer over all their backends. Useful for dashboards, sharing queries, and onboarding.

```bash
owit server --port 8080  # starts the API + serves the UI
```

---

## Architecture

![Architecture diagram](./public/openwatchit%20architecture.jpg)

**Core is written in Go.** Reasons: first-class gRPC support, excellent concurrency model for parallel dispatch, single binary distribution, strong CLI ecosystem (`cobra`, `viper`), and broad familiarity in the DevOps/Platform engineering community.

### Control Panel components

| Component | Responsibility |
|---|---|
| **OWL Parser** | Parses OWL query text into an AST |
| **Query Planner** | Converts the AST into an execution plan: which backends, join strategy, fan-out scope |
| **Plugin Manager** | Manages plugin process lifecycle; caches `Capabilities` responses; provides gRPC connections to the Dispatcher and Query Planner |
| **Dispatcher** | Sends the AST to target plugins in parallel; normalises incoming result rows against the canonical schema before forwarding |
| **Result Merger** | Combines normalised rows from all backends; executes cross-signal joins when no single backend covers all signal types |
| **API** | Exposes query execution to the CLI (Embedded and Remote Mode) and the Browser UI |

The **Renderer** is a CLI-layer component — it formats result rows and warnings for terminal output (table, JSON, stream). It is not part of the Control Panel. The Browser UI renders its own output independently.

---

## Vendor Support (Initial Targets)

| Vendor | Signals | Status |
|---|---|---|
| Grafana Loki | Logs | Planned v0.1 |
| Prometheus | Metrics | Planned v0.1 |
| Grafana Tempo | Traces | Planned v0.1 |
| Grafana Mimir | Metrics | Planned v0.1 |
| Jaeger | Traces | Planned v0.2 |
| Datadog | Logs, Metrics, Traces, Profiles | Planned v0.2 |
| AWS CloudWatch | Logs, Metrics | Planned v0.2 |
| Azure Application Insights | Logs, Metrics, Traces | Planned v0.3 |
| Google Cloud Logging | Logs | Community |
| New Relic | Logs, Metrics, Traces | Community |
| Splunk | Logs | Community |
| Elastic / OpenSearch | Logs | Community |
| Honeycomb | Traces, Logs | Community |

---

## Configuration

A single TOML file describes all your backends:

```toml
[backends.production-loki]
plugin  = "loki"
url     = "https://loki.internal"
default = true
tags    = ["prod", "logs"]

[backends.datadog]
plugin  = "datadog"
api_key = "$DD_API_KEY"
site    = "datadoghq.com"
default = true
tags    = ["prod", "logs", "metrics", "traces"]

[backends.staging-cloudwatch]
plugin  = "cloudwatch"
region  = "us-east-1"
profile = "staging"
default = false
tags    = ["staging", "logs", "metrics"]
```

You can query by tag:

```bash
owit logs --where "level == 'error'" --last 1h --tag prod
```

---

## Why Go. Why gRPC. Why Now.

**Go** gives us a single, self-contained binary. No runtime to install. `brew install openwatchit` and you're done. The concurrency primitives (goroutines, channels) are a natural fit for parallel backend queries.

**gRPC** for plugins means:
- Plugins are isolated processes — a crashing plugin doesn't crash core.
- Language-agnostic — write your plugin in whatever your team already knows.
- Strongly typed contracts — schema changes are explicit and versioned.
- Streaming support — query results can stream as they arrive.

**Now** because observability tooling is fragmenting faster than any single vendor can consolidate it. The Kubernetes ecosystem proved that the answer to fragmentation is open standards and plugin ecosystems, not winner-takes-all products.

---

## What We Are Not

- **Not a storage backend.** We don't store your data. We query it where it lives.
- **Not a replacement for your current tools.** We're a layer on top.
- **Not a paid SaaS.** OpenWatchIt is MIT-licensed. Run it anywhere.
- **Not opinionated about your stack.** If your team uses Datadog, great. If it uses a mix of six things, also great.

---

## Community Discussion

This is an early RFC. Nothing is final. Head to [**GitHub Discussions**](../../discussions) to share your perspective on the open questions that will shape what OpenWatchIt prioritizes.

We want to hear from DevOps engineers, SREs, software architects, and platform engineers.

- **Star this repo** if you'd use this tool today.
- **Share it** with the person on your team who would have the strongest opinion.
- **Open a discussion** with your own questions or use cases we haven't considered.

---

## Contributing

If this proposal gets enough signal, the next step is a working prototype of the core parser and two initial plugins (Loki and Prometheus). Watch this repo to follow along.

Read [**CONTRIBUTING.md**](CONTRIBUTING.md) to understand what meaningful contribution looks like at this stage.

---

_OpenWatchIt is not a company. It's a proposal. What it becomes depends entirely on whether the problem resonates with you._
