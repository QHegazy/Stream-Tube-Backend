package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomString(length int) string {
	randomBytes := make([]byte, (length*6+7)/8)
	rand.Read(randomBytes)
	encoded := base64.RawURLEncoding.EncodeToString(randomBytes)
	return encoded[:length]
}

func GenerateRecoveryCodes(count int) []string {
	codes := make([]string, count)
	for i := 0; i < count; i++ {
		codes[i] = GenerateRandomString(4)
	}
	return codes
}
