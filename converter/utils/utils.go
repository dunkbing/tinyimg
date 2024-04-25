package utils

import (
	"crypto/sha256"
	"fmt"
)

func GenerateHash(str string) (string, error) {
	hash := sha256.New()
	hash.Write([]byte(str))
	hashBytes := hash.Sum(nil)
	hashStr := fmt.Sprintf("%x", hashBytes)

	return hashStr, nil
}
