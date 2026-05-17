# ADR-0009: Backend credentials in TOML are env var references only (v0.x)

## Status
Accepted — revisit for v1.0 with secret manager support

## Context
The TOML config must reference backend credentials (API keys, tokens, passwords). A policy is needed on whether credential values can appear in plaintext in the config file.

## Decision
In v0.x, the TOML config only accepts env var references for sensitive fields. The syntax is `$VAR_NAME`:

```toml
[backends.datadog]
plugin  = "datadog"
api_key = "$DD_API_KEY"
```

The engine resolves `$VAR_NAME` at startup from the process environment. If the env var is not set, the engine fails with a clear error naming the missing variable. Plaintext credential values are rejected with a parse error.

## Reasons
- Establishes a secure-by-default culture from day one. No contributor can accidentally ship a config with hardcoded credentials.
- Env vars are the standard credential injection mechanism in CI/CD pipelines, Docker, and Kubernetes — the primary deployment environments for this tool.
- If the TOML config is committed to version control by mistake, no credentials are exposed.

## Consequences
- The TOML parser must distinguish between a regular string value and an env var reference (`$`-prefixed), and resolve the latter at startup.
- Error messages must clearly state which env var is missing and which backend requires it.
- **Future work (v1.0+):** extend the credential resolver to support external secret managers via a pluggable syntax (e.g. `vault:secret/path#key`, `awssm:secret-name`). The env var mechanism remains valid alongside these.
