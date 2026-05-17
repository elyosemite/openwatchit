# ADR-0008: Plugins install to user directory by default, local project directory with --local

## Status
Accepted

## Context
`owit plugin install <name>` downloads a plugin binary. A location on disk must be chosen. Options considered: system-global directory, user home directory, or project-local directory.

## Decision
**Default:** plugins install to `~/.owit/plugins/` (user home directory). No admin/root privileges required.

**Override:** `owit plugin install --local <name>` installs to `.owit/plugins/` relative to the directory containing the active TOML config. This allows pinning specific plugin versions per project and checking the plugin list into version control.

The CLI resolves plugins in this order: local `.owit/plugins/` first, then `~/.owit/plugins/`.

## Reasons
- User-directory install is the standard for developer CLI tools (analogous to `cargo install`, `npm -g` with `prefix` set to home). No permission escalation required.
- `--local` enables reproducible environments for teams: pin `tempo@2.1.0` in `.owit/plugins/` and everyone on the team uses the same version without coordinating a global install.
- Resolution order (local first) lets projects override a globally installed plugin version without affecting other projects.

## Consequences
- The CLI must check both locations at plugin spawn time and prefer local.
- `.owit/plugins/` should be added to `.gitignore` by default (binaries), but the team may choose to commit the directory for full reproducibility.
- The Plugin Marketplace must serve versioned, platform-specific binaries (at minimum: linux/amd64, linux/arm64, darwin/arm64, windows/amd64).
