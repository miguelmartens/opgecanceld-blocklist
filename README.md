# Opgecanceld

A domain blocklist focused on blocking YouTube advertisements, tracking, and ad-related services. Use this list to reduce or block YouTube ads across your devices and ad-blocking tools.

## Name

**Opgecanceld** is a Dutch neologism derived from _opgedonderd_ (a mild Dutch expletive) and _gecanceld_ (cancelled). A playful nod to cancel culture—here we cancel ads instead. Ads get _opgecanceld_.

Inspired by the show [AllStarsZonen](https://www.avrotros.nl/programmas/allstarszonen) / [Rundfunk](https://www.avrotros.nl/programmas/rundfunk). See [this clip](https://www.youtube.com/shorts/zd6ULJ1jLdo).

## Lists

| Format                  | Link                                                   | Compatible with                                                                                                        |
| ----------------------- | ------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------- |
| Domain list             | [opgecanceld-blocklist.txt](opgecanceld-blocklist.txt) | Pi-hole, AdGuard, AdGuard Home, uBlock Origin, AdAway, Blokada, DNS66, pfBlockerNG, Blocky, Technitium DNS, hosts file |
| AdGuard / uBlock filter | [opgecanceld-filters.txt](opgecanceld-filters.txt)     | AdGuard, uBlock Origin, AdBlock Plus                                                                                   |

## Usage

### Pi-hole / AdGuard Home / DNS-based blockers

Add the list URL to your blocklist sources:

```
https://raw.githubusercontent.com/<your-username>/opgecanceld-blocklist/main/opgecanceld-blocklist.txt
```

### Hosts file

Append the domains to your hosts file with `0.0.0.0` (or `127.0.0.1`):

- **Windows:** `C:\Windows\System32\drivers\etc\hosts`
- **macOS / Linux:** `/etc/hosts`

### Ad-blocker extensions (uBlock Origin, AdGuard, etc.)

**Domain list** (for Pi-hole-style DNS blocking):

```
https://raw.githubusercontent.com/<your-username>/opgecanceld-blocklist/main/opgecanceld-blocklist.txt
```

**AdGuard / uBlock filter list** (for browser extensions):

```
https://raw.githubusercontent.com/<your-username>/opgecanceld-blocklist/main/opgecanceld-filters.txt
```

Add the appropriate URL as a custom filter subscription in your ad-blocker.

## What this list blocks

- YouTube ad servers and tracking domains
- Google ad services (DoubleClick, Google Ads)
- Video ad networks (Innovid, Moat, etc.)
- Analytics and tracking scripts used for ad delivery

## Compatibility

- **Pi-hole** ✓
- **AdGuard / AdGuard Home** ✓
- **uBlock Origin** ✓
- **AdAway** (Android, rooted) ✓
- **Blokada** ✓
- **DNS66** ✓
- **pfBlockerNG** ✓
- **Hosts file** ✓

## Building

The AdGuard/uBlock filter list is generated from the domain blocklist. After editing `opgecanceld-blocklist.txt`, regenerate the filters:

```bash
make build-filters
# or
./bin/discover -build-filters
```

## Discovering new ad domains

The `discover` tool uses a headless browser to capture network traffic from YouTube and extract ad-related domains. It filters requests against known ad patterns (googlevideo, doubleclick, googlesyndication, etc.) and reports domains not yet in the blocklist.

**Requirements:** Chrome or Chromium must be installed (chromedp will find it automatically). Go 1.26+ to build from source.

```bash
# Build and run (captures for 2 minutes by default)
make run

# Or with go run
go run ./cmd/discover/

# Shorter capture for testing (30 seconds)
./bin/discover -duration 30s

# Save new domains to a file
./bin/discover -output new-domains.txt

# Append new domains directly to the blocklist (also regenerates filters)
./bin/discover -append
```

**Makefile targets:** `make build`, `make build-filters`, `make run`, `make dev`, `make test`, `make lint`, `make fmt`, `make install`, `make format`, `make format-check`, `make lint-yaml`, `make renovate`

`make run` and `make dev` append new domains and regenerate the filter list automatically.

## Development

- **Go:** `make fmt` (format code, tidy modules), `make test`, `make lint` (golangci-lint)
- **Markdown/YAML/JSON:** `make format-check` (CI check) or `make format` / `make prettier` to fix. Prettier is pinned to 3.3.2.
- **YAML:** `make lint-yaml` (yamllint)

## Automated Discovery

A [scheduled workflow](.github/workflows/discover-scheduled.yml) runs weekly (Mondays at 6am UTC) to discover new YouTube ad domains and open a pull request with any updates:

1. Runs the discover tool with a 30-second capture
2. Appends new domains to the blocklist and regenerates filters
3. Creates a PR with the changes (or updates an existing PR)
4. Auto-approves and enables auto-merge so the PR merges without manual intervention

Trigger manually via **Actions → Discover ad domains → Run workflow**.

**Required repo settings** for auto-merge to work:

- **Settings → General → Pull Requests:** Enable "Allow auto-merge"
- **Settings → Branches:** Add a branch protection rule for `main` with at least one requirement (e.g. "Require status checks to pass")

## Automated Dependency Management

This project uses [Renovate](https://docs.renovatebot.com/) for automated dependency updates:

- Automatic PRs for Go modules and GitHub Actions
- Scheduled weekly updates (Mondays before 6am UTC)

**Setup:** Install the [Renovate GitHub App](https://github.com/apps/renovate) on the repository. Merge the onboarding PR Renovate creates.

## License

MIT License — see [LICENSE](LICENSE) for details.
