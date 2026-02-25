package blocklist

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDomainSet(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "blocklist.txt")

	// Write a sample blocklist
	content := `# Comment line
example.com
doubleclick.net
# Another comment

googlevideo.com
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	got, err := LoadDomainSet(path)
	if err != nil {
		t.Fatalf("LoadDomainSet: %v", err)
	}

	want := map[string]bool{
		"example.com":     true,
		"doubleclick.net": true,
		"googlevideo.com": true,
	}
	for d, ok := range want {
		if !got[d] {
			t.Errorf("expected %q in set", d)
		}
		_ = ok
	}
	if len(got) != len(want) {
		t.Errorf("got %d domains, want %d", len(got), len(want))
	}
}

func TestLoadDomainSet_Nonexistent(t *testing.T) {
	got, err := LoadDomainSet("/nonexistent/path/blocklist.txt")
	if err != nil {
		t.Fatalf("LoadDomainSet on nonexistent file should not error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty set for nonexistent file, got %d", len(got))
	}
}
