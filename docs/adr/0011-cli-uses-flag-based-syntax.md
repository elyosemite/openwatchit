# ADR-0011: OWL CLI syntax uses positional filter triplets and named flags

## Status
Accepted

## Context
Early designs wrapped OWL queries in a quoted string (`owit query "traces | where duration > 1s | ..."`), which conflicted with the shell `|` operator and required escaping. A flag-based approach (`--where "level == 'error'"`) improved on this but still required quoting for each filter expression. A pipeline syntax for files or interactive use was also considered and dropped entirely.

Additionally, OWL was initially described as "KQL-inspired." This framing was dropped: OWL is its own language.

## Decision
OWL is expressed exclusively as a CLI command. The signal type is the subcommand; filters are positional triplets; options are named flags.

**Filter triplets** — `field op value`, multiple triplets space-separated:
```bash
owit logs level eq error service eq checkout --last 30m --limit 100
owit traces duration gt 500ms root_error eq true --last 1h --limit 20
owit logs level ne info message contains timeout --last 1h --limit 50
```

**Equality shorthand** — `field=value` is syntactic sugar for `field eq value`:
```bash
owit logs level=error service=checkout --last 30m --limit 100
```

**Filter operators** (OData-inspired vocabulary — `eq`, `ne`, `gt`, `ge`, `lt`, `le`, `contains`).

**Named flags** for non-filter options: `--last`, `--limit`, `--summarize`, `--backends`, `--tag`, `--output`.

The CLI parser collects positional arguments as filter triplets until it encounters the first `--flag`, which avoids ambiguity without a separator token.

There is no pipeline syntax and no file-based query format in v0.x.

## Reasons
- Positional filter triplets require no quoting for any common filter — no shell conflicts at all.
- The `field=value` shorthand keeps equality filters concise without sacrificing readability.
- OData operator names (`eq`, `ne`, `gt`, `ge`, `lt`, `le`) are already familiar to developers who use Azure CLI, REST APIs, or OData-based tools — zero learning curve for the vocabulary.
- A single, consistent surface (CLI flags) removes the cognitive overhead of choosing between syntaxes.
- Fully discoverable: `owit logs --help` lists all available operators and flags without any external documentation.

## Consequences
- The CLI parser must implement triplet grouping for positional args and shorthand `field=value` expansion to `field eq value`.
- Each signal type subcommand shares the same filter syntax — no signal-type-specific parsing rules.
- New filter operators added in future versions require parser changes but no flag additions.
- The `--summarize` flag accepts a quoted expression (e.g. `"count() by service"`) — the only flag that still requires quoting due to the space in its value.
