---
title: Installation
description: How to install the owit CLI.
---

## Requirements

- Linux (amd64 / arm64) or macOS (arm64)
- No runtime dependencies — `owit` is a single self-contained binary

## Install via Homebrew

```bash
brew install openwatchit
```

## Install via script

```bash
curl -fsSL https://install.openwatchit.dev | sh
```

## Download manually

Download the latest binary for your platform from the [GitHub Releases](https://github.com/elyosemite/openwatchit/releases) page.

```bash
# Example for Linux amd64
curl -L https://github.com/elyosemite/openwatchit/releases/latest/download/owit_linux_amd64 \
  -o /usr/local/bin/owit && chmod +x /usr/local/bin/owit
```

## Verify the installation

```bash
owit --version
```

## Next steps

- [Configure your backends](/getting-started/configuration/)
