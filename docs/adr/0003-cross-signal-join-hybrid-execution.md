# ADR-0003: Cross-signal joins use hybrid execution (plugin-native when possible, core fallback)

## Status
Accepted

## Context
OWL supports joining across signal types (e.g. `logs | join traces on traceId`). This requires querying two different signal types, potentially from different vendor backends, and correlating the results.

Three approaches were considered:
- **A**: Core always executes the join in the Result Merger after two parallel fan-outs.
- **B**: The plugin executes the join natively when the vendor supports multiple signal types.
- **C**: Hybrid — route to a single plugin when one backend covers all signal types involved; fall back to core-side join otherwise.

## Decision
Hybrid execution (Option C):

1. The Query Planner inspects the join and identifies all signal types involved.
2. It queries the `Capabilities` of registered backends to find any single backend that declares support for all signal types in the join.
3. If such a backend exists (e.g. DataDog supporting both `logs` and `traces`), the full join AST is routed to that plugin for native execution.
4. If no single backend covers all signal types, the core executes two parallel fan-outs and performs the join in the Result Merger using the Normalized Result schema.

## Reasons
- The cross-vendor join (core path) is the product's key differentiator — it must work regardless of vendor overlap.
- Plugin-native joins are a performance optimization for cases like DataDog (which has native correlation between logs and traces). Forcing the core path there wastes a capability the vendor already has.
- The Query Planner already queries Capabilities for fan-out routing; the extra check for join coverage adds minimal complexity.

## Consequences
- The `Capabilities` gRPC response must declare supported signal types per backend explicitly.
- Plugin authors for multi-signal vendors (DataDog, New Relic, Honeycomb) should implement native join translation to take advantage of the plugin path.
- The Result Merger must implement a general in-memory join on the Normalized Result schema for the core fallback path.
- In v0.x, only the core fallback path is required. Plugin-native join translation is optional and can be added incrementally.
