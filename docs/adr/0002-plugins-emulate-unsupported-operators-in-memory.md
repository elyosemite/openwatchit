# ADR-0002: Plugins emulate unsupported OWL operators in memory (v0.x)

## Status
Accepted — revisit before v1.0

## Context
Some OWL operators have no direct equivalent in a vendor's native query language. Example: `limit N` in OWL has no counterpart in PromQL (the closest is `topk(k, <vector>)`, which is semantically different — it ranks by value, not by arrival order).

Requiring every plugin to perform deep semantic mapping between OWL operators and vendor-native equivalents in v0.x is too large a scope for initial development.

## Decision
In v0.x, plugins may emulate unsupported operators by fetching the full result set and applying the operation in memory within the plugin process before returning the Normalized Result stream to the engine.

## Reasons
- Unblocks initial plugin development without requiring deep vendor-language expertise for every operator.
- Keeps the plugin contract simple: receive AST, return normalized rows.
- The correctness of results is preserved; only efficiency is sacrificed.

## Consequences
- **Known cost**: plugins may transfer significantly more data than necessary over the wire (e.g. fetching 100k log lines to return 50).
- **Future work**: before v1.0, plugins should map common OWL operators to semantically equivalent vendor-native constructs where possible (e.g. `limit` → `topk` for ranked Prometheus queries, `limit` → Loki's native line limit). This mapping can be incremental per operator per vendor.
- Plugin authors should document which operators are natively translated vs. emulated in their plugin's README.
