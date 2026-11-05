// Package peinfo extracts static metadata from Windows PE files.
package peinfo

import (
	"bytes"
	"debug/pe"
	"fmt"
)

type Info struct {
	Machine  uint16   `json:"machine"`
	Sections []string `json:"sections"`
	Imports  []string `json:"imports"`
}

func Parse(data []byte) (*Info, error) {
	f, err := pe.NewFile(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("parse PE: %w", err)
	}
	defer f.Close()
	result := &Info{Machine: f.FileHeader.Machine}
	for _, section := range f.Sections {
		result.Sections = append(result.Sections, section.Name)
	}
	if imports, err := f.ImportedLibraries(); err == nil {
		result.Imports = imports
	}
	return result, nil
}
