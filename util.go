package main

import (
	"encoding/base64"
	"fmt"
	"regexp"

	"golang.org/x/mod/sumdb/dirhash"
)

func calculateDirChecksum(path string, prefix string) (string, error) {
	// calculate sha256 checksum, which is returned as a base64 encoded string
	// prefixed with "h1:"
	h1, err := dirhash.HashDir(path, prefix, dirhash.Hash1)
	if err != nil {
		return "", fmt.Errorf("error calculating checksum of dir: %w", err)
	}

	// remove the "h1:" prefix
	re := regexp.MustCompile(`^h1:`)
	b64 := re.ReplaceAllString(h1, "")

	// base64 decode string
	bytes, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", fmt.Errorf("error base64 decoding checksum: %w", err)
	}

	return fmt.Sprintf("%x", bytes), nil
}
