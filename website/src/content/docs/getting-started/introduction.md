---
title: Introduction
description: What OpenWatchIt is and why it exists.
---

OpenWatchIt is an open-source observability gateway. It lets you query **logs, metrics, traces, and profiles** across any combination of vendors using a single unified command-line tool.

## The problem

Modern engineering teams accumulate observability tools organically. Logs in Loki. Metrics in Prometheus. Traces in DataDog. Each tool has its own query language, its own UI, and its own mental model.

During an incident, this means:

- Opening four browser tabs
- Remembering four query syntaxes
- Manually copying trace IDs between tools
- Writing one-off scripts to correlate data

The tools are individually excellent. **The seams between them are the problem.**

## The solution

OpenWatchIt sits in front of your existing tools. You write one query; it dispatches to every configured backend in parallel, translates to each vendor's native language, and returns a single normalised result.

```bash
owit logs level=error service=payments --last 1h --limit 50
```

That single command hits Loki, DataDog, and CloudWatch simultaneously — one result, sorted by timestamp, with a `_source` column showing where each row came from.

## What OpenWatchIt is not

- **Not a storage backend.** Your data stays where it is.
- **Not a replacement for your current tools.** It's a layer on top.
- **Not a paid SaaS.** MIT-licensed. Run it anywhere.

## Next steps

- [Install OpenWatchIt](/getting-started/installation/)
- [Configure your backends](/getting-started/configuration/)
