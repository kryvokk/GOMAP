// Package hasher provides cryptographic hashes for files.
package hasher

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// FileSHA256 returns the SHA256 hash of the file at path as a lowercase
// hexadecimal string.
func FileSHA256(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
