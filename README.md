# GoMalScan

GoMalScan is a command-line application for static analysis of potentially malicious files. It is designed to inspect files without executing them.

## Current status

The project currently provides the initial CLI foundation: help and version commands. File scanning and analysis will be added in later development phases.

## Requirements

- Go 1.27.3 or a compatible newer Go toolchain

## Build and run

```sh
go build -o gomal-scan ./cmd/gomal-scan
./gomal-scan --help
./gomal-scan version
```

## Project layout

- `cmd/gomal-scan` contains the executable entry point.
- `internal/cli` contains command-line behavior.
- `docs` contains architecture and development notes.

## Security

GoMalScan is intended to analyze untrusted files without executing them. The current CLI does not yet scan files. Never treat heuristic findings from future versions as definitive proof of malware.
