package smhb

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

// Get data path prefix
func prefix(context serverContext, extension string) string {
	return context.dataPath + "/" + extension
}

// Hash password
func hash(cleartext string) []byte {
	hash := sha256.New()
	hash.Write([]byte(cleartext))
	return hash.Sum(nil)
}

// Serialize generic data using gob
func serialize(this interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	if err := encoder.Encode(this); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// Deserialize generic data using gob
func deserialize(this interface{}, buffer []byte) error {
	reader := bytes.NewReader(buffer)
	decoder := gob.NewDecoder(reader)

	if err := decoder.Decode(this); err != nil {
		return err
	}

	return nil
}
