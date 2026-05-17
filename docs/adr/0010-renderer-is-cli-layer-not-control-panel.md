# ADR-0010: The Renderer is a CLI-layer component, not part of the Control Panel

## Status
Accepted

## Context
The original component list included a Renderer inside the Control Panel alongside the OWL Parser, Query Planner, Fan-out Executor, Result Merger, and API.

## Decision
The Renderer is removed from the Control Panel and classified as a CLI-layer component.

The Control Panel produces structured data only: Normalized Result rows and a warnings array. Formatting that data for human consumption (table, JSON, stream) is the responsibility of the CLI, which calls the Renderer after receiving the result — directly in Embedded Mode, or after deserializing the API response in Remote Mode.

## Reasons
- The Browser UI (React) connects to the API and renders results with its own components. If the Renderer lived in the Control Panel, it would be a component that the primary UI client (the browser) never uses.
- The API boundary already implies structured data exchange. Adding a rendering step inside the Control Panel would mean the server is formatting terminal output — a category error.
- Different CLI commands may render the same result differently (`--output table` vs `--output json` vs `owit tail` streaming). This variation belongs at the CLI layer, not inside the engine.

## Consequences
- The Control Panel API response schema must be self-contained and format-agnostic: rows, warnings, metadata.
- The CLI owns all output formatting decisions: column widths, color, streaming behavior, truncation.
- The Browser UI is fully decoupled from CLI rendering — it can evolve its own display independently.
