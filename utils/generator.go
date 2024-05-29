package utils

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

// Function to generate a random string
func RandString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
