# GoMalScan project plan

## Phase 1 — Foundation

- [x] Go module and executable entry point
- [x] CLI help and version commands
- [x] Basic CLI tests
- [x] Initial project documentation
- [ ] CI pipeline

## Phase 2 — File scanner

- [x] Read individual files as a stream and calculate SHA256
- [ ] Read directories safely
- Handle invalid paths, permissions, symbolic links, and size limits
- Identify file types

## Later phases

Signature detection, metadata analysis, static PE analysis, explainable risk scoring, reports, local history, and optional integrations will follow once earlier phases are stable. Scanned files must never be executed or uploaded automatically.
