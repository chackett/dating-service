package security

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func CreateSecureSessionToken(n int) (string, error) {
	token := make([]byte, n)

	_, err := rand.Read(token)
	if err != nil {
		return "", fmt.Errorf("failed to generate secure token: %w", err)
	}

	return hex.EncodeToString(token), nil
}
