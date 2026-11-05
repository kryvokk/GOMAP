package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/kryvokk/GOMAP/internal/scanner"
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
	root.AddCommand(&cobra.Command{
		Use:   "scan <file-or-directory>...",
		Short: "Scan files and directories without executing files",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fileScanner := scanner.Scanner{}
			for _, path := range args {
				results, err := fileScanner.ScanPath(path)
				if err != nil {
					return err
				}
				for _, result := range results {
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "SHA256  %s  SIZE  %d  TYPE  %s  %s\n", result.SHA256, result.Size, result.MIMEType, result.Path); err != nil {
						return err
					}
				}
			}
			return nil
		},
	})
	return root
}

// Execute runs the CLI using the process arguments and standard output.
func Execute() error {
	return newRootCommand(os.Stdout).Execute()
}
