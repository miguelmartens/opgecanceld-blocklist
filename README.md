# Opgecanceld

A domain blocklist focused on blocking YouTube advertisements, tracking, and ad-related services. Use this list to reduce or block YouTube ads across your devices and ad-blocking tools.

## Name

**Opgecanceld** is a Dutch neologism derived from *opgedonderd* (a mild Dutch expletive) and *gecanceld* (cancelled). A playful nod to cancel culture—here we cancel ads instead. Ads get *opgecanceld*.

Inspired by the show [AllStarsZonen](https://www.avrotros.nl/programmas/allstarszonen) / [Rundfunk](https://www.avrotros.nl/programmas/rundfunk). See [this clip](https://www.youtube.com/shorts/zd6ULJ1jLdo).

## Lists

| Format | Link | Compatible with |
|--------|------|-----------------|
| Domain list | [opgecanceld-blocklist.txt](opgecanceld-blocklist.txt) | Pi-hole, AdGuard, AdGuard Home, uBlock Origin, AdAway, Blokada, DNS66, pfBlockerNG, Blocky, Technitium DNS, hosts file |
| AdGuard / uBlock filter | [opgecanceld-filters.txt](opgecanceld-filters.txt) | AdGuard, uBlock Origin, AdBlock Plus |

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
python3 build-filters.py
```

## License

MIT License — see [LICENSE](LICENSE) for details.
