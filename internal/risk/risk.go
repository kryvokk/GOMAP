// Package risk assigns explainable heuristic risk levels. Findings are indicators,
// not proof that a file is malicious.
package risk

type Finding struct {
	Rule   string `json:"rule"`
	Points int    `json:"points"`
	Reason string `json:"reason"`
}
type Assessment struct {
	Score    int       `json:"score"`
	Level    string    `json:"level"`
	Findings []Finding `json:"findings"`
}

func Assess(signatureNames []string, hasMZ, peParsed bool) Assessment {
	var a Assessment
	for _, name := range signatureNames {
		a.Findings = append(a.Findings, Finding{"local-signature", 70, "Matched local byte signature: " + name})
	}
	if hasMZ && !peParsed {
		a.Findings = append(a.Findings, Finding{"invalid-pe", 25, "MZ marker is present but PE metadata could not be parsed"})
	}
	for _, f := range a.Findings {
		a.Score += f.Points
	}
	if a.Score > 100 {
		a.Score = 100
	}
	switch {
	case a.Score >= 70:
		a.Level = "high"
	case a.Score >= 25:
		a.Level = "medium"
	default:
		a.Level = "low"
	}
	return a
}
