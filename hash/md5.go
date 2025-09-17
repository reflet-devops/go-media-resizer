package hash

import (
	"crypto/md5"
	"encoding/hex"
)

/*func GenerateMD5(content io.Reader) (string, error) {
	hasher := md5.New()
	_, err := io.Copy(hasher, content)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}*/

func GenerateMD5FromString(content string) (string, error) {
	return GenerateMD5FromBytes([]byte(content))
}

func GenerateMD5FromBytes(content []byte) (string, error) {
	hasher := md5.New()
	hasher.Write(content)
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
