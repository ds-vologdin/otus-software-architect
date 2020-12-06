package service

import (
	"crypto/sha256"
	"encoding/base64"
)

// HashString get hash of string
func HashString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return base64.StdEncoding.EncodeToString(hash[:])
}
