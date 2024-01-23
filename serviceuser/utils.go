package serviceuser

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

func HashAPIKey(apiKey string) string {
	hashedBytes := sha256.Sum256([]byte(apiKey))
	return hex.EncodeToString(hashedBytes[:])
}

func GenerateAPIKey() (string, error) {
	apiKey, err := generateRandomString(32)
	if err != nil {
		return "", err
	}

	return apiKey, nil
}

func generateRandomString(length int) (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	for i := range bytes {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		bytes[i] = chars[n.Int64()]
	}
	return string(bytes), nil
}
