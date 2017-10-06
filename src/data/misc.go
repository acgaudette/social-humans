package data

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

// Get data path prefix
func prefix(extension string) string {
	return DATA_PATH + "/" + extension
}

// Hash password
func hash(cleartext string) []byte {
	hash := sha256.New()
	hash.Write([]byte(cleartext))
	return hash.Sum(nil)
}

// Generate random string
func generateToken() string {
	buffer := make([]byte, 32)
	rand.Read(buffer)
	return fmt.Sprintf("%x", buffer)
}
