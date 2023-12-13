// Generates Ethereum compatible private key, public key, and address
// Writes keys and address to files
// Private key format: 64 character hex string with '0x' prefix
// Public key format: 128 character uncompressed hex string with "0x" prefix, not including "0x04" prefix
// Usage: go run scripts/gen-keys/main.go

package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	// Generate a new private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	// Get the associated public key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	// Serialize the public key in uncompressed format
	pubKeyBytes := append(padByteSlice(publicKeyECDSA.X.Bytes()), padByteSlice(publicKeyECDSA.Y.Bytes())...)

	// Get the address
	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Convert keys to strings with '0x' prefix
	privateKeyBytes := fmt.Sprintf("0x%x", crypto.FromECDSA(privateKey))
	publicKeyString := fmt.Sprintf("0x%x", pubKeyBytes)

	// Write private key to file
	err = os.WriteFile("private_key.txt", []byte(privateKeyBytes), 0644)
	if err != nil {
		log.Fatalf("Failed to write private key to file: %v", err)
	}

	// Write public key to file
	err = os.WriteFile("public_key.txt", []byte(publicKeyString), 0644)
	if err != nil {
		log.Fatalf("Failed to write public key to file: %v", err)
	}

	// Write address to file
	err = os.WriteFile("address.txt", []byte(address.Hex()), 0644)
	if err != nil {
		log.Fatalf("Failed to write address to file: %v", err)
	}

	fmt.Println("Keys and address have been written to files.")
}

// padByteSlice pads a byte slice to 32 bytes
func padByteSlice(b []byte) []byte {
	padded := make([]byte, 32)
	copy(padded[32-len(b):], b)
	return padded
}
