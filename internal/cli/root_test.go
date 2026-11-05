package cli

import (
	"bytes"
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
