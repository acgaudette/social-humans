package smhb

import (
	"io/ioutil"
	"log"
)

// Return addresses of all post files for a given user handle
func getPostAddresses(author string, context ServerContext) ([]string, error) {
	// Read posts directory for user
	files, err := ioutil.ReadDir(prefix(context, author+"/"))

	if err != nil {
		return nil, err
	}

	addresses := []string{}

	// Build addresses slice
	for _, file := range files {
		// Get address from filename
		addresses = append(addresses, author+"/"+file.Name()[0:len(file.Name())-5])
	}

	log.Printf("Collated post addresses for \"%s\"", author)

	return addresses, nil
}

// Serialize post addresses for a given user handle
func serializePostAddresses(
	author string, context ServerContext,
) ([]byte, error) {
	addresses, err := getPostAddresses(author, context)

	if err != nil {
		return nil, err
	}

	return serialize(addresses)
}

// Deserialize raw buffer with lookup handle
func deserializePostAddresses(buffer []byte) ([]string, error) {
	addresses := []string{}
	err := deserialize(&addresses, buffer)

	if err != nil {
		return nil, err
	}

	return addresses, nil
}
