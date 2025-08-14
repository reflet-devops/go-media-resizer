package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"strings"
)

func GenerateSHA256(content io.Reader) (string, error) {
	hasher := sha256.New()
	_, err := io.Copy(hasher, content)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func GenerateSHA256FromString(content string) (string, error) {
	return GenerateSHA256(strings.NewReader(content))
}
