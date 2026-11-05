package scanner

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestScanPathFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sample.txt")
	if err := os.WriteFile(path, []byte("hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	results, err := (Scanner{}).ScanPath(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("got %d results, want 1", len(results))
	}
	if results[0].SHA256 != "5891b5b522d5df086d0ff0b110fbd9d21bb4fc7163af34d08286a2e846f6be03" {
		t.Errorf("unexpected hash: %s", results[0].SHA256)
	}
	if results[0].Size != 6 {
		t.Errorf("size = %d, want 6", results[0].Size)
	}
	if results[0].MIMEType == "" {
		t.Error("MIME type is empty")
	}
}

func TestScanPathDirectory(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"a.txt", "b.txt"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(name), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	results, err := (Scanner{}).ScanPath(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
}

func TestScanPathMissing(t *testing.T) {
	_, err := (Scanner{}).ScanPath(filepath.Join(t.TempDir(), "missing"))
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestScanPathRejectsOversizedFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "large.bin")
	if err := os.WriteFile(path, []byte("12345"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := (Scanner{MaxFileSize: 4}).ScanPath(path)
	if !errors.Is(err, ErrFileTooLarge) {
		t.Fatalf("error = %v, want ErrFileTooLarge", err)
	}
}

func TestScanPathRejectsSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target.txt")
	if err := os.WriteFile(target, []byte("data"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "link.txt")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	_, err := (Scanner{}).ScanPath(link)
	if err == nil {
		t.Fatal("expected symlink to be rejected")
	}
}
