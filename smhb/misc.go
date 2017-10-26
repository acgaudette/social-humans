package smhb

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
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

func serialize(this interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	if err := encoder.Encode(this); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func deserialize(this interface{}, buffer []byte) error {
	reader := bytes.NewReader(buffer)
	decoder := gob.NewDecoder(reader)

	if err := decoder.Decode(this); err != nil {
		return err
	}

	return nil
}
