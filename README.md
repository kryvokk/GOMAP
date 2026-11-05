# GoMalScan

GoMalScan is a local command-line tool for static inspection of files. It never launches scanned files. Findings and risk levels are heuristic indicators and do not prove that a file is infected.

## Build

Requires Go 1.27.3 or a compatible newer toolchain.

```sh
go build -o gomal-scan ./cmd/gomal-scan
./gomal-scan --help
```

## Quick start

```sh
./gomal-scan scan ./suspicious-file.exe
./gomal-scan scan ./folder --format json --output report.json
./gomal-scan scan ./folder --format html --output report.html
./gomal-scan history
```

The scanner reports SHA-256, size, MIME type, extension, modification time, MZ/PE markers, PE headers, sections and imported libraries. The default per-file read limit is 64 MiB. Symbolic links and non-regular files (including devices and FIFOs) are reported and skipped. Errors for individual files do not stop the rest of a directory scan.

## Local signatures

Create a UTF-8 text file with one named hex byte sequence per line. Blank lines and lines beginning with `#` are ignored.

```text
# name=hex bytes
demo-marker=4d5a
test-string=746573742d7061796c6f6164
```

Use `--signatures rules.txt`. A matched local rule is an indicator only. The scanner reads file contents only up to the configured limit.

## Reports and options

`scan` accepts one or more paths. Options:

- `--format text|json|html` selects the report format.
- `--output FILE` writes the report to a file.
- `--signatures FILE` loads local byte signatures.
- `--yara-rules FILE` optionally invokes the external `yara` executable; set `--yara-bin` if it is elsewhere on `PATH`.
- `--ti-url URL` optionally performs a hash-only GET to `URL/files/{sha256}`. `--ti-key` sends a bearer token. This generic endpoint contract may require an adapter for a particular provider. Hash lookups disclose file hashes and should only be enabled with a trusted service.

Risk scores are rule-based and explainable. A high score or signature match is not a malware verdict. PE parsing is static and does not execute code.

## Scan history

Completed scans are appended locally to `gomalscan/history.jsonl` under the operating system's user config directory (permissions are restricted to the user where supported). `gomal-scan history` prints the 20 latest records. The JSON Lines layout is intentionally simple; a later migration can move it to SQLite.

## Docker

```sh
docker build -t gomalscan .
docker run --rm -v "$PWD:/scan:ro" gomalscan scan /scan --format json
```

YARA is not installed in the image. Mount or install a YARA executable and its rules separately if that optional integration is needed. The container is non-root and the input mount should be read-only.

## Development

```sh
go test ./...
go build ./cmd/gomal-scan
```
