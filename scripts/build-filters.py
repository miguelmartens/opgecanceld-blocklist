#!/usr/bin/env python3
"""
Generate AdGuard / uBlock Origin filter list from opgecanceld-blocklist.txt.
Run from project root: python3 scripts/build-filters.py
"""

import os

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
ROOT_DIR = os.path.dirname(SCRIPT_DIR)

BLOCKLIST = os.path.join(ROOT_DIR, "opgecanceld-blocklist.txt")
OUTPUT = os.path.join(ROOT_DIR, "opgecanceld-filters.txt")

HEADER = """! --------------------------------------------
! Opgecanceld - AdGuard / uBlock Origin filter list
! --------------------------------------------
! Title: Opgecanceld
! Description: Blocks YouTube ads, tracking, and ad-related services.
! Homepage: https://github.com/your-username/opgecanceld-blocklist
! License: MIT
! Expires: 4 days
!
"""


def main():
    with open(BLOCKLIST) as f:
        lines = f.readlines()

    domains = set()
    for line in lines:
        line = line.strip()
        if not line or line.startswith("#"):
            continue
        domains.add(line)

    rules = [f"||{d}^" for d in sorted(domains)]
    output = HEADER + "\n".join(rules) + "\n"

    with open(OUTPUT, "w") as f:
        f.write(output)

    print(f"Generated {OUTPUT} with {len(rules)} filter rules")


if __name__ == "__main__":
    main()
