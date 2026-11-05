package main

import (
	"os"

	"github.com/kryvokk/GOMAP/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
