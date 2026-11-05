package cli

import (
	"context"
	"fmt"
	"github.com/kryvokk/GOMAP/internal/history"
	"github.com/kryvokk/GOMAP/internal/report"
	"github.com/kryvokk/GOMAP/internal/signatures"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

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
			yaraRules, _ := cmd.Flags().GetString("yara-rules")
			yaraBin, _ := cmd.Flags().GetString("yara-bin")
			tiURL, _ := cmd.Flags().GetString("ti-url")
			tiKey, _ := cmd.Flags().GetString("ti-key")
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
			for i := range all {
				if all[i].Error != "" {
					continue
				}
				if yaraRules != "" {
					ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Second)
					out, runErr := exec.CommandContext(ctx, yaraBin, yaraRules, all[i].Path).CombinedOutput()
					cancel()
					if runErr != nil && len(out) == 0 {
						return fmt.Errorf("run YARA on %q: %w", all[i].Path, runErr)
					}
					for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
						if fields := strings.Fields(line); len(fields) > 0 {
							all[i].YARAMatches = append(all[i].YARAMatches, fields[0])
						}
					}
				}
				if tiURL != "" {
					u, err := url.Parse(strings.TrimRight(tiURL, "/") + "/files/" + all[i].SHA256)
					if err != nil {
						return err
					}
					ctx, cancel := context.WithTimeout(cmd.Context(), 10*time.Second)
					req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
					if err != nil {
						cancel()
						return err
					}
					if tiKey != "" {
						req.Header.Set("Authorization", "Bearer "+tiKey)
					}
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						cancel()
						return fmt.Errorf("threat-intelligence lookup: %w", err)
					}
					resp.Body.Close()
					cancel()
					if resp.StatusCode < 200 || resp.StatusCode >= 300 {
						return fmt.Errorf("threat-intelligence lookup returned %s", resp.Status)
					}
					all[i].ThreatIntel = fmt.Sprintf("HTTP %d", resp.StatusCode)
				}
			}
			historyPath, err := history.DefaultPath()
			if err != nil {
				return err
			}
			if err := history.Append(historyPath, all); err != nil {
				return fmt.Errorf("save scan history: %w", err)
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
	root.AddCommand(&cobra.Command{Use: "history", Short: "Show recent local scans", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, args []string) error {
		path, err := history.DefaultPath()
		if err != nil {
			return err
		}
		records, err := history.Read(path, 20)
		if err != nil {
			return err
		}
		for _, record := range records {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s  %d file entries\n", record.Time.Format(time.RFC3339), len(record.Results)); err != nil {
				return err
			}
		}
		return nil
	}})
	scan := root.Commands()[1]
	scan.Flags().String("format", "text", "report format: text, json, or html")
	scan.Flags().String("output", "", "write report to this file")
	scan.Flags().String("signatures", "", "load local signature rules (name=hex-bytes per line)")
	scan.Flags().String("yara-rules", "", "optionally run external yara CLI rules against scanned files")
	scan.Flags().String("yara-bin", "yara", "path to YARA executable")
	scan.Flags().String("ti-url", "", "optional threat-intelligence API base URL (hash-only GET /files/{sha256})")
	scan.Flags().String("ti-key", "", "optional bearer token for threat-intelligence API")
	return root
}

// Execute runs the CLI using the process arguments and standard output.
func Execute() error {
	return newRootCommand(os.Stdout).Execute()
}
