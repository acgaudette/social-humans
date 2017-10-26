package smhb

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

// Get data path prefix
func prefix(extension string) string {
	return DATA_PATH + "/" + extension
}

// The address is a unique string identifier for the post
func BuildPostAddress(handle, stamp string) string {
	return handle + "/" + stamp
}

// Hash password
func hash(cleartext string) []byte {
	hash := sha256.New()
	hash.Write([]byte(cleartext))
	return hash.Sum(nil)
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
