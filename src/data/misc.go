package data

import (
	"io/ioutil"
	"strings"
)

// Read the short commit hash from file
func getCommitHash(root string) (*string, error) {
	// Read HEAD
	buffer, err := ioutil.ReadFile(root + "/HEAD")

	if err != nil {
		return nil, err
	}

	// Convert to string and trim newline
	path := string(buffer[:])
	path = path[:len(path)-1]

	// Get path from HEAD ref
	ref := strings.Split(path, " ")[1]

	// Read commit hash
	buffer, err = ioutil.ReadFile(root + "/" + ref)

	if err != nil {
		return nil, err
	}

	// Convert to string
	hash := string(buffer[:])

	// Make short hash
	hash = hash[:7]

	return &hash, nil
}
