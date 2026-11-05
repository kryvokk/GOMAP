package hasher

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileSHA256(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sample.txt")
	if err := os.WriteFile(path, []byte("hello\n"), 0o600); err != nil {
		t.Fatalf("write sample file: %v", err)
	}

	got, err := FileSHA256(path)
	if err != nil {
		t.Fatalf("FileSHA256() error = %v", err)
	}
	const want = "5891b5b522d5df086d0ff0b110fbd9d21bb4fc7163af34d08286a2e846f6be03"
	if got != want {
		t.Fatalf("FileSHA256() = %q, want %q", got, want)
	}
}

func TestFileSHA256MissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.txt")

	got, err := FileSHA256(path)
	if err == nil {
		t.Fatal("FileSHA256() error = nil, want an error")
	}
	if got != "" {
		t.Fatalf("FileSHA256() hash = %q, want an empty string on error", got)
	}
}
