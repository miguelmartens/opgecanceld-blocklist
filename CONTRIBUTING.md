# Contributing to Opgecanceld Blocklist

Thank you for your interest in contributing. This document explains how to get set up, run checks, and open a pull request.

## Requirements

- **To run the discover tool:** **macOS**, **Linux**, or **Windows** (Chrome or Chromium required).
- **To build, test, and lint:** install locally:
  - **Go 1.26+** — [Install Go](https://go.dev/dl/)
  - Optional for local checks: **golangci-lint**, **prettier**, **yamllint** (CI will run these on your PR if you don't have them).

## Getting started

1. **Fork** the repository on GitHub, then clone your fork:

   ```bash
   git clone https://github.com/YOUR_USERNAME/opgecanceld-blocklist.git
   cd opgecanceld-blocklist
   ```

2. **Build and run** to confirm everything works:

   ```bash
   make build
   make run
   ```

   The discover tool will capture ad domains from YouTube. See [README.md](README.md) for usage.

## Making changes

1. **Create a branch** for your change:

   ```bash
   git checkout -b feat/my-feature
   ```

2. **Make your edits.** The project follows [Standard Go Project Layout](https://github.com/golang-standards/project-layout):

   - `cmd/discover/` — main entry point
   - `internal/blocklist/` — blocklist loading and filter generation
   - `internal/discover/` — YouTube ad domain discovery (chromedp)

3. **Run checks** before opening a PR:

   ```bash
   make test
   make lint
   make format-check
   make lint-yaml
   ```

   Fix any failures so CI passes.

4. **Commit and push** to your fork, then open a **pull request** with a short description of the change.

## Development commands

| Command              | Description                                        |
| -------------------- | -------------------------------------------------- |
| `make build`         | Build binary to `bin/discover`                     |
| `make build-filters` | Generate AdGuard/uBlock filter list from blocklist |
| `make run`           | Build and run (discovers and appends new domains)  |
| `make dev`           | Clean, then build and run                          |
| `make test`          | Run tests                                          |
| `make lint`          | Run golangci-lint                                  |
| `make fmt`           | Format Go code and tidy modules                    |
| `make format`        | Format Markdown/YAML/JSON with Prettier            |
| `make format-check`  | Check formatting (CI)                              |
| `make lint-yaml`     | Run yamllint on YAML files                         |

## Pull request process

- Keep PRs focused; prefer smaller changes over large ones.
- For bugs or new features, opening an **issue** first is welcome (not required).
- CI runs on every push and PR: **build** (build + test), **lint**, **format-check**, **lint-yaml**. All must pass before merge.
- Maintainers will review and may request changes.

## Further reading

- [README.md](README.md) — Usage, blocklist formats, discovery tool
- [docs/RENOVATE_SETUP.md](docs/RENOVATE_SETUP.md) — Renovate configuration
