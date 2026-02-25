// Package blocklist provides loading and writing of domain blocklist files.
package blocklist

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// LoadDomainSet reads a blocklist file and returns a set of lowercase domain names.
// Comments (lines starting with #) and empty lines are ignored.
// If the file does not exist, returns an empty set.
func LoadDomainSet(path string) (map[string]bool, error) {
	existing := make(map[string]bool)
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return existing, nil
		}
		return nil, fmt.Errorf("open blocklist: %w", err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		existing[strings.ToLower(line)] = true
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("read blocklist: %w", err)
	}
	return existing, nil
}

// AppendDomains appends domains to a blocklist file with a section header.
func AppendDomains(path string, domains []string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open blocklist: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString("\n# --- Discovered domains ---\n"); err != nil {
		return fmt.Errorf("write header: %w", err)
	}
	for _, d := range domains {
		if _, err := fmt.Fprintln(f, d); err != nil {
			return fmt.Errorf("write domain: %w", err)
		}
	}
	return nil
}

// WriteDomains writes domains to a file, one per line.
func WriteDomains(path string, domains []string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer f.Close()

	for _, d := range domains {
		if _, err := fmt.Fprintln(f, d); err != nil {
			return fmt.Errorf("write domain: %w", err)
		}
	}
	return nil
}

// AdGuardFilterHeader is the header for the generated AdGuard/uBlock filter list.
const AdGuardFilterHeader = `! --------------------------------------------
! Opgecanceld - AdGuard / uBlock Origin filter list
! --------------------------------------------
! Title: Opgecanceld
! Description: Blocks YouTube ads, tracking, and ad-related services.
! Homepage: https://github.com/your-username/opgecanceld-blocklist
! License: MIT
! Expires: 4 days
!
`

// GenerateFilters reads the blocklist and writes the AdGuard/uBlock filter list.
func GenerateFilters(blocklistPath, outputPath string) (int, error) {
	domains, err := LoadDomainSet(blocklistPath)
	if err != nil {
		return 0, fmt.Errorf("load blocklist: %w", err)
	}

	var sorted []string
	for d := range domains {
		sorted = append(sorted, d)
	}
	sort.Strings(sorted)

	var b strings.Builder
	b.WriteString(AdGuardFilterHeader)
	for _, d := range sorted {
		b.WriteString("||")
		b.WriteString(d)
		b.WriteString("^\n")
	}

	if err := os.WriteFile(outputPath, []byte(b.String()), 0o644); err != nil {
		return 0, fmt.Errorf("write filters: %w", err)
	}
	return len(sorted), nil
}

