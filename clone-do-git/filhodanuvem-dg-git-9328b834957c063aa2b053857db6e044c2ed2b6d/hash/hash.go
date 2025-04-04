package hash

import (
	"crypto/sha1"
	"encoding/hex"
)

func New(fileContent []byte) string {
	sha := sha1.New()
	sha.Write(fileContent)

	return hex.EncodeToString(sha.Sum(nil))
}
