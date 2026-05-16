# OpenWatchIt

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

You write one query. OpenWatchIt fans it out to every configured backend in parallel, translates it to each vendor's native language under the hood, and returns a unified, normalized result.

```
owit query "logs | where level == 'error' and service == 'payments' | last 1h | limit 50"
```

That single line hits Loki, Datadog, and CloudWatch simultaneously — you get one result, ranked by timestamp, with a `_source` column showing where each entry came from.

No new vendor to learn. No new dashboard to configure. One tool, every source.

---

## Core Pillars

### 1. OpenWatch Query Language (OWL)

A KQL-inspired, pipeline-based query language that maps to every signal type:

```sql
-- Logs
logs | where level == "error" | where service == "checkout" | last 30m | limit 100

-- Metrics
metrics | where __name__ == "http_requests_total" | where env == "prod"
       | summarize sum(value) by service | order by sum desc

-- Traces
traces | where duration > 500ms | where root_error == true
       | project traceId, service, duration, spanCount | limit 20

-- Profiles
profiles | where type == "cpu" | where service == "api-gateway"
         | summarize avg(value) by function | order by avg desc | limit 10

-- Cross-signal correlation (the real power)
logs | where level == "error"
    | join traces on traceId
    | project timestamp, message, duration, spanId
```

OWL is intentionally readable. A developer who has never used it before should be able to write a useful query in under five minutes.

### 2. Multi-vendor Fan-out

Queries are executed against all configured — or explicitly targeted — backends in parallel:

```bash
# All configured backends
owit query "logs | where level == 'error'"

# Specific backends
owit query --backends loki,datadog "logs | where service == 'api'"

# All backends of a given type
owit query --type logs "logs | where level == 'error'"
```

Results arrive as backends respond. Slow backends don't block fast ones.

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
owit query "traces | where duration > 1s" --output table
owit query "metrics | where ..." --output json | jq '.[] | .value'
owit repl  # interactive mode with autocomplete
owit tail "logs | where service == 'api'"  # streaming, like tail -f
```

**UI** — a browser-based interface that connects to an OpenWatchIt server instance, designed for teams who want a shared, visual layer over all their backends. Useful for dashboards, sharing queries, and onboarding.

```bash
owit server --port 8080  # starts the API + serves the UI
```

---

## Architecture

```
┌──────────────────────────────────────────────────┐
│              CLI  /  Browser UI                  │
└──────────────────┬───────────────────────────────┘
                   │
┌──────────────────▼───────────────────────────────┐
│               OpenWatchIt Core  (Go)             │
│                                                  │
│  OWL Parser → Query Planner → Fan-out Executor   │
│  Result Merger → Normalizer → Renderer / API     │
└──────┬──────────┬──────────┬──────────┬──────────┘
       │ gRPC     │ gRPC     │ gRPC     │ gRPC
┌──────▼──┐ ┌────▼────┐ ┌───▼────┐ ┌──▼──────────┐
│  Loki   │ │DataDog  │ │CloudW. │ │  Any Plugin  │
│ Plugin  │ │ Plugin  │ │ Plugin │ │  (community) │
│  (Go)   │ │  (Go)   │ │  (Go)  │ │  (any lang)  │
└─────────┘ └─────────┘ └────────┘ └─────────────┘
```

**Core is written in Go.** Reasons: first-class gRPC support, excellent concurrency model for fan-out, single binary distribution, strong CLI ecosystem (`cobra`, `viper`), and broad familiarity in the DevOps/Platform engineering community.

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
owit query --tag prod "logs | where level == 'error' | last 1h"
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

---

_OpenWatchIt is not a company. It's a proposal. What it becomes depends entirely on whether the problem resonates with you._
