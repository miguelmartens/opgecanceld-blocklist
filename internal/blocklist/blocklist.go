// Package blocklist provides loading and writing of domain blocklist files.
package blocklist

import (
	"bufio"
	"fmt"
	"os"
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
