package auth

import (
	"crypto/rand"
	"fmt"
)

const (
	TokenLength = 128
	charset     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func GenerateToken(length int) (string, error) {
	charsetLen := len(charset)
	token := make([]byte, length)

	// Fill token with random bytes from crypto/rand
	_, err := rand.Read(token)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	// Convert bytes to characters from charset
	for i := 0; i < length; i++ {
		token[i] = charset[int(token[i])%charsetLen]
	}

	return string(token), nil
}
