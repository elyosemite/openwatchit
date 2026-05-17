# ADR-0005: CLI embeds the Control Panel for solo use; optionally connects to a remote server for team use

## Status
Accepted

## Context
The README describes two usage patterns: `owit query` (direct CLI use by an engineer) and `owit server` (a running server that also serves the Browser UI). A decision was needed on whether the Control Panel is always a separate process or can be embedded in the CLI binary.

## Decision
The `owit` binary operates in two modes:

**Embedded mode (default):** `owit query`, `owit tail`, `owit repl` spin up the Control Panel in-process. No server required. Plugins are spawned as child processes by the CLI and terminated when the command finishes. The user's TOML config is read from the local filesystem.

**Remote mode:** If a `server_url` is set in the TOML config (or passed via `--server`), the CLI acts as a thin client — it sends the OWL query to the remote Control Panel over HTTP/gRPC and streams back results. The remote server manages plugins and backend config centrally.

`owit server --port 8080` starts the Control Panel as a long-running process, exposes the API, and serves the Browser UI.

## Reasons
- A solo engineer can install `owit`, drop a TOML config, and run queries immediately — no infrastructure to stand up.
- A team can run one `owit server` instance with centrally managed backend credentials and config. Engineers point their CLI at it and share the same backend pool without each having local credentials.
- The Browser UI requires a running server regardless; this model makes `owit server` a natural fit for team setups.

## Consequences
- The Control Panel logic must be usable both as an in-process library (for embedded mode) and as a network service (for server mode). This suggests a clean internal interface boundary between the engine core and its transport layer (CLI invocation vs. HTTP/gRPC server).
- Plugin processes in embedded mode are short-lived (spawned and killed per command). In server mode, plugins are long-running child processes of the server. Plugin authors should design for both lifecycles.
- Local TOML config in embedded mode means backend credentials live on the engineer's machine. Remote mode centralises credential management on the server. Teams should prefer remote mode for credential hygiene.
