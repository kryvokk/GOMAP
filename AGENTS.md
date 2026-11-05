# GoMalScan contributor instructions

- Keep the application static-analysis-only: never execute analyzed files or upload them automatically.
- Treat all scanned files and parsed content as untrusted input.
- Prefer the Go standard library; explain any new dependency and keep dependencies focused.
- Keep packages small and responsibilities clear. Application implementation belongs under `internal/`.
- Use English for code comments and `gofmt` for Go source.
- Add tests for new behavior and update relevant documentation.
- Do not commit or push changes without explicit approval.
