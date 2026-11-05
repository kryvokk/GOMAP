package cli

import (
	"fmt"
	"github.com/kryvokk/GOMAP/internal/report"
	"github.com/kryvokk/GOMAP/internal/signatures"
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
			format, _ := cmd.Flags().GetString("format")
			rulePath, _ := cmd.Flags().GetString("signatures")
			outputPath, _ := cmd.Flags().GetString("output")
			fileScanner := scanner.Scanner{}
			if rulePath != "" {
				rules, err := signatures.Load(rulePath)
				if err != nil {
					return err
				}
				fileScanner.Signatures = rules
			}
			var all []scanner.Result
			for _, path := range args {
				results, err := fileScanner.ScanPath(path)
				if err != nil {
					return err
				}
				all = append(all, results...)
			}
			w := cmd.OutOrStdout()
			if outputPath != "" {
				f, err := os.Create(outputPath)
				if err != nil {
					return err
				}
				defer f.Close()
				w = f
			}
			return report.Write(w, format, all)
		},
	})
	scan := root.Commands()[1]
	scan.Flags().String("format", "text", "report format: text, json, or html")
	scan.Flags().String("output", "", "write report to this file")
	scan.Flags().String("signatures", "", "load local signature rules (name=hex-bytes per line)")
	return root
}

// Execute runs the CLI using the process arguments and standard output.
func Execute() error {
	return newRootCommand(os.Stdout).Execute()
}
