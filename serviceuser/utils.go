package serviceuser

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateAPIKey() (string, error) {

	// Create a byte slice of 32 bytes
	bytes := make([]byte, 32)

	// Read random bytes, using the crypto/rand package for cryptographic security
	if _, err := rand.Read(bytes); err != nil {
		// Handle the error appropriately in your real code
		return "", err
	}

	// Encode the bytes to a base64 string
	return base64.URLEncoding.EncodeToString(bytes), nil
}
