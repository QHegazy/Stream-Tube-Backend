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