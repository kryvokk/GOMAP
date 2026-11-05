# Architecture

The executable lives in `cmd/gomal-scan`. It delegates command behavior to the private `internal/cli` package. Future scanning and analysis packages will be introduced only when their responsibilities are implemented, and the dependency flow should remain one-directional toward reusable internal components.
