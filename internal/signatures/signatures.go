// Package signatures loads and matches local byte-sequence rules.
package signatures

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

type Rule struct {
	Name  string
	Bytes []byte
}

func Load(path string) ([]Rule, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open signature rules: %w", err)
	}
	defer f.Close()
	var rules []Rule
	s := bufio.NewScanner(f)
	for line := 1; s.Scan(); line++ {
		value := strings.TrimSpace(s.Text())
		if value == "" || strings.HasPrefix(value, "#") {
			continue
		}
		parts := strings.SplitN(value, "=", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" {
			return nil, fmt.Errorf("rule line %d: expected name=hex-bytes", line)
		}
		b, err := hex.DecodeString(strings.ReplaceAll(strings.TrimSpace(parts[1]), " ", ""))
		if err != nil || len(b) == 0 {
			return nil, fmt.Errorf("rule line %d: invalid hex byte sequence", line)
		}
		rules = append(rules, Rule{Name: strings.TrimSpace(parts[0]), Bytes: b})
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("read signature rules: %w", err)
	}
	return rules, nil
}

func Match(data []byte, rules []Rule) []string {
	var found []string
	for _, rule := range rules {
		if len(rule.Bytes) > 0 && strings.Contains(string(data), string(rule.Bytes)) {
			found = append(found, rule.Name)
		}
	}
	return found
}
