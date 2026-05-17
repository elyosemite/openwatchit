# ADR-0001: Plugins receive a parsed AST, not raw OWL text

## Status
Accepted

## Context
The plugin gRPC contract includes a `Translate` RPC. A plugin must convert an OWL query into its vendor's native query language (PromQL, LogQL, DogStatsD, KQL, etc.). Two approaches were considered:

- **Option A**: pass raw OWL text — the plugin parses OWL itself and emits native syntax.
- **Option B**: the core parses OWL once, producing an AST, and passes the AST to the plugin.

## Decision
Plugins receive a protobuf-encoded AST, not raw OWL text.

## Reasons
1. The OWL parser is the most complex, specialized component in the project. Duplicating it in every plugin language (Go, Rust, Java, Python…) is infeasible.
2. The AST is a stable, versioned protobuf schema. OWL syntax can evolve without breaking plugin contracts — plugins consume semantic structure, not surface syntax.
3. Centralising parsing enables consistent error messages, query validation, and future language features across all plugins without any plugin changes.

## Consequences
- The core must expose a well-defined, versioned AST schema in the `.proto` files.
- Plugin authors never write an OWL parser — they only walk an AST and emit their vendor's native query.
- Breaking changes to the AST schema require a major version bump and a migration path for existing plugins.
