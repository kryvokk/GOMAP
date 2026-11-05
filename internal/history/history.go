// Package history stores completed scans locally as append-only JSON Lines.
package history

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kryvokk/GOMAP/internal/scanner"
)

type Record struct {
	Time    time.Time        `json:"time"`
	Results []scanner.Result `json:"results"`
}

func DefaultPath() (string, error) {
	d, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "gomalscan", "history.jsonl"), nil
}
func Append(path string, results []scanner.Result) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(Record{Time: time.Now().UTC(), Results: results})
}
func Read(path string, limit int) ([]Record, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var records []Record
	s := bufio.NewScanner(f)
	s.Buffer(make([]byte, 64*1024), 32*1024*1024)
	for s.Scan() {
		var r Record
		if err := json.Unmarshal(s.Bytes(), &r); err != nil {
			return nil, fmt.Errorf("decode history: %w", err)
		}
		records = append(records, r)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if limit > 0 && len(records) > limit {
		records = records[len(records)-limit:]
	}
	return records, nil
}
