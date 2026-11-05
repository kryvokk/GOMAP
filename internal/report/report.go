package report

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/kryvokk/GOMAP/internal/scanner"
)

func Write(w io.Writer, format string, results []scanner.Result) error {
	switch strings.ToLower(format) {
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	case "html":
		const page = `<!doctype html><html><head><meta charset="utf-8"><title>GoMalScan report</title><style>body{font:14px sans-serif;margin:2rem}table{border-collapse:collapse;width:100%}th,td{border:1px solid #ccc;padding:.5rem;text-align:left}th{background:#eee}.high{color:#b00}</style></head><body><h1>GoMalScan report</h1><p>Heuristic indicators do not prove malware infection.</p><table><tr><th>Path</th><th>SHA256</th><th>Size</th><th>Type</th><th>Risk</th><th>Findings / Error</th></tr>{{range .}}<tr><td>{{.Path}}</td><td>{{.SHA256}}</td><td>{{.Size}}</td><td>{{.MIMEType}}</td><td class="{{.Risk.Level}}">{{.Risk.Level}} ({{.Risk.Score}})</td><td>{{if .Error}}{{.Error}}{{else}}{{range .Risk.Findings}}{{.Reason}}; {{end}}{{range .Signatures}}signature: {{.}}; {{end}}{{end}}</td></tr>{{end}}</table></body></html>`
		t, err := template.New("report").Parse(page)
		if err != nil {
			return err
		}
		return t.Execute(w, results)
	case "text", "terminal", "":
		for _, r := range results {
			if r.Error != "" {
				if _, err := fmt.Fprintf(w, "ERROR  %s  %s\n", r.Path, r.Error); err != nil {
					return err
				}
				continue
			}
			if _, err := fmt.Fprintf(w, "SHA256  %s  SIZE  %d  TYPE  %s  RISK  %s(%d)  %s\n", r.SHA256, r.Size, r.MIMEType, r.Risk.Level, r.Risk.Score, r.Path); err != nil {
				return err
			}
			for _, f := range r.Risk.Findings {
				if _, err := fmt.Fprintf(w, "  [%d] %s\n", f.Points, f.Reason); err != nil {
					return err
				}
			}
		}
		return nil
	default:
		return fmt.Errorf("unsupported report format %q (use text, json, or html)", format)
	}
}
