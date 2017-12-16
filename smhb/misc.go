package smhb

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

// Get data path prefix
func prefix(context ServerContext, extension string) string {
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

	// Encode struct with gob
	if err := encoder.Encode(this); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// Deserialize generic data using gob
func deserialize(this interface{}, buffer []byte) error {
	reader := bytes.NewReader(buffer)
	decoder := gob.NewDecoder(reader)

	// Decode buffer with gob
	if err := decoder.Decode(this); err != nil {
		return err
	}

	return nil
}

// Serialize and prepend a string timestamp to serialized data
func prependTimestamp(data []byte, timestamp string) ([]byte, error) {
	time, err := serialize(timestamp)
	if err != nil {
		return nil, err
	}
	return append(time, data...), nil
}

// Returns serialized data without timestamp, the timestamp, and error
func extractTimestamp(request REQUEST, data []byte) ([]byte, string, error) {
	var timestamp *string
	reader := bytes.NewReader(data)
	decoder := gob.NewDecoder(reader)

	if err := decoder.Decode(data); err != nil {
		return nil, "", err
	}

	buffer := &bytes.Buffer{}
	buffer.ReadFrom(reader)
	data = buffer.Bytes()

	return data, *timestamp, nil
}
