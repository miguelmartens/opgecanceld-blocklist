# Renovate Setup Guide

Automated dependency management for Opgecanceld Blocklist.

## Quick Start

1. Go to [Renovate GitHub App](https://github.com/apps/renovate)
2. Click **Configure**
3. Select this repository (or your fork)
4. Grant the requested permissions
5. Click **Install**
6. Merge the onboarding PR Renovate creates
7. Check the Dependency Dashboard issue for available updates

## What Renovate Does

- Opens PRs for outdated **Go modules** (`go.mod` / `go.sum`)
- Opens PRs for outdated **GitHub Actions** in `.github/workflows/`
- Runs on the configured schedule (see below)
- Labels PRs with `dependencies`

## Configuration

The [renovate.json](../renovate.json) at the repo root configures:

| Feature     | Setting                |
| ----------- | ---------------------- |
| Base preset | `config:recommended`   |
| Schedule    | Mondays before 6am UTC |
| Go modules  | Enabled                |

## Local Dry-Run

To preview what Renovate would do without opening PRs:

```bash
make renovate
```

## Resources

- [Renovate Docs](https://docs.renovatebot.com/)
- [Config Validator](https://app.renovatebot.com/config-validator)
