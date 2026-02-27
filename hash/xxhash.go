package hash

import (
	"encoding/hex"

	"github.com/cespare/xxhash/v2"
)

func GenerateXXHashFromString(content string) (string, error) {
	return GenerateXXHashFromBytes([]byte(content))
}

func GenerateXXHashFromBytes(content []byte) (string, error) {
	hasher := xxhash.New()
	_, _ = hasher.Write(content)
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
