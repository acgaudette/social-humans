package smhb

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"strings"
)

// Return addresses of all post files for a given user handle
func getPostAddresses(context serverContext, author string) ([]string, error) {
	// Read posts directory for user
	files, err := ioutil.ReadDir(prefix(context, author+"/"))

	if err != nil {
		return nil, err
	}

	addresses := []string{}

	// Build addresses slice
	for _, file := range files {
		// Get address from filename
		address := author + "/" + strings.Split(file.Name(), ".")[0]
		addresses = append(addresses, address)
	}

	return addresses, nil
}

// Serialize post addresses for a given user handle
func serializePostAddresses(
	context serverContext, author string,
) ([]byte, error) {
	addresses, err := getPostAddresses(context, author)

	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	if err = encoder.Encode(addresses); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// Deserialize raw buffer with lookup handle
func deserializePostAddresses(buffer []byte) ([]string, error) {
	addresses := []string{}

	reader := bytes.NewReader(buffer)
	decoder := gob.NewDecoder(reader)

	if err := decoder.Decode(&addresses); err != nil {
		return nil, err
	}

	return addresses, nil
}
