# GoMalScan project plan

## Phase 1 — Foundation

- [x] Go module and executable entry point
- [x] CLI help and version commands
- [x] Basic CLI tests
- [x] Initial project documentation
- [ ] CI pipeline

## Phase 2 — File scanner

- Read individual files and directories safely
- Handle invalid paths, permissions, symbolic links, and size limits
- Calculate SHA256 and identify file types

## Later phases

Signature detection, metadata analysis, static PE analysis, explainable risk scoring, reports, local history, and optional integrations will follow once earlier phases are stable. Scanned files must never be executed or uploaded automatically.
