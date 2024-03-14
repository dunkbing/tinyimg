package utils

import (
	"crypto/sha256"
	"fmt"
)

func GenerateHash(buf []byte) (string, error) {
	hash := sha256.New()
	hash.Write(buf[:32*1024])
	hashBytes := hash.Sum(nil)
	hashStr := fmt.Sprintf("%x", hashBytes)

	return hashStr, nil
}
