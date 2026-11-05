// Package scanner reads files as untrusted data and collects basic metadata.
package scanner

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const DefaultMaxFileSize int64 = 64 << 20

var ErrFileTooLarge = errors.New("file exceeds maximum allowed size")

type Result struct {
	Path     string
	Size     int64
	SHA256   string
	MIMEType string
	Error    string
}

type Scanner struct {
	MaxFileSize int64
}

func (s Scanner) ScanPath(path string) ([]Result, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("inspect %q: %w", path, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("scan %q: symbolic links are not supported", path)
	}
	if info.IsDir() {
		var results []Result
		err := filepath.WalkDir(path, func(current string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				results = append(results, Result{Path: current, Error: walkErr.Error()})
				return nil
			}
			if entry.Type()&os.ModeSymlink != 0 {
				results = append(results, Result{Path: current, Error: "symbolic links are not supported"})
				return nil
			}
			if entry.IsDir() {
				return nil
			}
			result, err := s.scanFile(current)
			if err != nil {
				results = append(results, Result{Path: current, Error: err.Error()})
				return nil
			}
			results = append(results, result)
			return nil
		})
		if err != nil {
			return results, fmt.Errorf("scan directory %q: %w", path, err)
		}
		return results, nil
	}
	result, err := s.scanFile(path)
	if err != nil {
		return nil, err
	}
	return []Result{result}, nil
}

func (s Scanner) scanFile(path string) (Result, error) {
	limit := s.MaxFileSize
	if limit <= 0 {
		limit = DefaultMaxFileSize
	}
	// Inspect the path before opening it. Opening a FIFO can otherwise block
	// indefinitely, and opening a device may have side effects.
	pathInfo, err := os.Lstat(path)
	if err != nil {
		return Result{}, fmt.Errorf("inspect %q: %w", path, err)
	}
	if pathInfo.Mode()&os.ModeSymlink != 0 {
		return Result{}, fmt.Errorf("scan %q: symbolic links are not supported", path)
	}
	if !pathInfo.Mode().IsRegular() {
		return Result{}, fmt.Errorf("scan %q: only regular files are supported", path)
	}
	file, err := os.Open(path)
	if err != nil {
		return Result{}, fmt.Errorf("open %q: %w", path, err)
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return Result{}, fmt.Errorf("stat %q: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return Result{}, fmt.Errorf("scan %q: only regular files are supported", path)
	}
	if info.Size() > limit {
		return Result{}, fmt.Errorf("scan %q: %w (%d bytes; limit %d)", path, ErrFileTooLarge, info.Size(), limit)
	}
	data, err := io.ReadAll(io.LimitReader(file, limit+1))
	if err != nil {
		return Result{}, fmt.Errorf("read %q: %w", path, err)
	}
	if int64(len(data)) > limit {
		return Result{}, fmt.Errorf("scan %q: %w (limit %d)", path, ErrFileTooLarge, limit)
	}
	digest := sha256.Sum256(data)
	mimeType := "application/octet-stream"
	if len(data) > 0 {
		mimeType = http.DetectContentType(data)
	}
	return Result{Path: path, Size: int64(len(data)), SHA256: hex.EncodeToString(digest[:]), MIMEType: mimeType}, nil
}
