package main

import (
	"fmt"
	"os"

	"github.com/kryvokk/GOMAP/internal/cli"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println(`   ██████╗  ██████╗ ███╗   ███╗ █████╗ ██████╗
  ██╔════╝ ██╔═══██╗████╗ ████║██╔══██╗██╔══██╗
  ██║  ███╗██║   ██║██╔████╔██║███████║██████╔╝
  ██║   ██║██║   ██║██║╚██╔╝██║██╔══██╗██╔═══╝
  ╚██████╔╝╚██████╔╝██║ ╚═╝ ██║██║  ██║██║
   ╚═════╝  ╚═════╝ ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝`)
		fmt.Println("  Static malware analysis · use --help for commands")
	}
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "gomal-scan:", err)
		os.Exit(1)
	}
}
