package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRootCommandShowsHelp(t *testing.T) {
	var output bytes.Buffer
	command := newRootCommand(&output)
	command.SetArgs([]string{"--help"})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(output.String(), "GoMalScan inspects files") {
		t.Fatalf("help output missing application description: %q", output.String())
	}
}

func TestScanCommand(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sample.txt")
	if err := os.WriteFile(path, []byte("hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	command := newRootCommand(&output)
	command.SetArgs([]string{"scan", path})
	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(output.String(), "5891b5b522d5df086d0ff0b110fbd9d21bb4fc7163af34d08286a2e846f6be03") {
		t.Fatalf("scan output missing expected hash: %q", output.String())
	}
}

func TestVersionCommand(t *testing.T) {
	originalVersion := Version
	Version = "1.2.3"
	t.Cleanup(func() { Version = originalVersion })

	var output bytes.Buffer
	command := newRootCommand(&output)
	command.SetArgs([]string{"version"})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if got, want := output.String(), "gomal-scan 1.2.3\n"; got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}
