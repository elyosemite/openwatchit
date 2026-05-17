# ADR-0006: Plugins in Embedded Mode run as lazy daemons with inactivity timeout

## Status
Accepted

## Context
In Embedded Mode, plugins are separate gRPC processes. The question is when these processes are started and when they are terminated.

Three options were considered: spawn-per-invocation, persistent daemon, or lazy daemon with inactivity timeout.

## Decision
Plugins in Embedded Mode are started on first use (lazy) and kept alive until an inactivity timeout elapses with no queries routed to them. The timeout is configurable in the TOML config (default: 30s). On timeout, the plugin process exits cleanly.

The CLI manages plugin PIDs via a local socket or PID file. On each invocation, it checks whether the required plugins are already running and reuses them; if not, it spawns them.

## Reasons
- Eliminates plugin startup latency for interactive use (e.g. multiple queries during an incident response session).
- Avoids leaving idle plugin processes running indefinitely on developer machines.
- The configurable timeout lets users tune for their workflow: lower for resource-constrained machines, higher for heavy interactive use.

## Consequences
- The CLI must implement a lightweight plugin process manager: spawn, health-check, reuse, and timeout-based shutdown.
- Plugin authors must ensure their process exits cleanly when the gRPC connection is closed or a shutdown signal is received.
- In Remote Mode (server), plugins are long-running and this timeout logic does not apply — the server manages plugin lifecycle independently.
