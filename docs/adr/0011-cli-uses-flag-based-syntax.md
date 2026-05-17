# ADR-0011: OWL is expressed exclusively through CLI flag-based syntax

## Status
Accepted

## Context
Early designs considered two OWL surfaces: a flag-based CLI syntax and a pipeline syntax (using `|` as operator separator) for files or interactive use. The pipeline syntax was dropped entirely. OWL has one surface only.

Additionally, OWL was initially described as "KQL-inspired." This framing was dropped: OWL is its own language and should not be marketed by reference to another product.

## Decision
OWL is expressed exclusively as CLI flags. The signal type is the subcommand; each operator is a named flag:

```bash
owit traces --where "duration > 1s" --where "root_error == true" --last 1h --limit 20
owit logs --where "level == 'error'" --last 30m --backends loki,datadog
owit tail logs --where "service == 'payments'"
```

There is no pipeline syntax, no `.owl` file format, and no `owit query "..."` string form.

## Reasons
- A single surface eliminates the need for a CLI adapter layer and a string parser — the CLI flags map directly to the AST.
- Flag-based syntax is natively discoverable: `owit traces --help` lists every available operator without any documentation lookup.
- Removing the pipeline syntax removes an entire category of shell-escaping bugs and user confusion about which syntax to use where.
- Dropping the KQL framing removes an implicit constraint on OWL's evolution — the language can develop its own identity.

## Consequences
- Each OWL operator (`where`, `last`, `limit`, `summarize`, and all future additions) must have a corresponding CLI flag.
- The OWL Parser receives a structured input (flags parsed by the CLI framework) rather than a string — the "parser" is effectively the flag-to-AST mapping layer in the CLI.
- There is no file-based query format in v0.x. Reusable queries are shell scripts or shell aliases.
