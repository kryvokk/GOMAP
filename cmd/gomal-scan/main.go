package main

import (
	"fmt"
	"os"

	"github.com/kryvokk/GOMAP/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "gomal-scan:", err)
		os.Exit(1)
	}
}
