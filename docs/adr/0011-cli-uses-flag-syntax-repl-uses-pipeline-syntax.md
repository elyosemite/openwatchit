# ADR-0011: CLI uses flag-based syntax; REPL uses OWL pipeline syntax

## Status
Accepted

## Context
OWL's pipeline syntax uses `|` as the operator separator. The shell also uses `|` as a pipe operator. Wrapping an entire OWL query in a quoted string (`owit query "traces | where duration > 1s | ..."`) avoids the conflict but feels artificial and creates friction — especially with string literals inside the query that require escaping.

Additionally, OWL was initially described as "KQL-inspired." This framing was dropped: OWL is its own language and should not be marketed by reference to another product.

## Decision
OWL has two surface syntaxes that both produce the same AST:

**CLI syntax** — used in the terminal. The signal type is the subcommand; each OWL operator maps to a named flag:
```bash
owit traces --where "duration > 1s" --where "root_error == true" --last 1h --limit 20
owit logs --where "level == 'error'" --last 30m --backends loki,datadog
owit tail logs --where "service == 'payments'"
```

**Pipeline syntax** — used in `owit repl` and `.owl` files:
```
traces | where duration > 1s | where root_error == true | last 1h | limit 20
```

A CLI adapter layer converts flag-based input into an OWL pipeline string before passing it to the OWL Parser. The Parser, the AST, and everything downstream are unaware of which surface was used.

## Reasons
- Flag-based CLI syntax is natural to terminal users: discoverable via `--help`, composable, no shell quoting of the full expression.
- Keeping a pipeline syntax for the REPL preserves expressiveness for users who want to write multi-operator queries interactively without flag ergonomics.
- A single parser processing a single canonical syntax keeps the engine simple.
- Dropping the KQL framing removes an implicit constraint on OWL's evolution — the language can develop its own identity.

## Consequences
- The CLI must implement a flag-to-OWL adapter for each signal type subcommand.
- Each OWL operator in v0.1 (`where`, `last`, `summarize`, `limit`) must have a corresponding CLI flag. New operators added in later versions must also get CLI flag equivalents.
- The `owit query "..."` form (raw OWL string) may be retained as a power-user/scripting escape hatch but is not the primary CLI UX.
- Documentation and examples must use CLI syntax for terminal use cases and pipeline syntax for REPL/file use cases — never mixing the two in the same context.
