---
title: Configuration
description: How to configure backends in the TOML config file.
---

OpenWatchIt is configured via a single TOML file. By default it is read from `~/.owit/config.toml`.

## Minimal example

```toml
[backends.production-loki]
plugin  = "loki"
url     = "https://loki.internal"
default = true
tags    = ["prod", "logs"]

[backends.prometheus]
plugin  = "prometheus"
url     = "https://prometheus.internal"
default = true
tags    = ["prod", "metrics"]
```

## Backend fields

| Field | Required | Description |
|---|---|---|
| `plugin` | yes | Plugin name to use for this backend |
| `url` | yes | Base URL of the vendor API |
| `default` | yes | Include in fan-out when no `--backends` flag is passed |
| `tags` | no | Arbitrary strings for targeting with `--tag` |

## Credentials

Sensitive values must be referenced as environment variables — never stored as plaintext.

```toml
[backends.datadog]
plugin  = "datadog"
api_key = "$DD_API_KEY"
app_key = "$DD_APP_KEY"
site    = "datadoghq.com"
default = true
tags    = ["prod", "logs", "metrics", "traces"]
```

If a referenced environment variable is not set, `owit` will fail at startup with a clear error naming the missing variable.

## Plugin timeout

Configure the inactivity timeout for plugin processes (default: 30s):

```toml
[plugins]
inactivity_timeout = "60s"
```

## Remote server

To connect to a shared `owit server` instead of running the engine locally:

```toml
server_url = "https://owit.internal:8080"
```

When `server_url` is set, the CLI acts as a thin client and all plugin management is handled by the server.
