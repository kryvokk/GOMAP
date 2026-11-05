package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// Version is set at build time with -ldflags when creating a release.
var Version = "dev"

func newRootCommand(out io.Writer) *cobra.Command {
	root := &cobra.Command{
		Use:           "gomal-scan",
		Short:         "Static analysis for potentially malicious files",
		Long:          "GoMalScan inspects files without executing them.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.SetOut(out)
	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the application version",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "gomal-scan %s\n", Version)
			return err
		},
	})
	return root
}

// Execute runs the CLI using the process arguments and standard output.
func Execute() error {
	return newRootCommand(os.Stdout).Execute()
}
