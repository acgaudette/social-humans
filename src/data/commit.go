package data

import (
	"io/ioutil"
	"log"
	"strings"
)

var CommitHash *string

// Load commit hash on initialization
func init() {
	hash, err := getCommitHash()

	if err != nil {
		log.Printf("%s", err)
		return
	}

	CommitHash = hash
}

func getCommitHash() (*string, error) {
	// Read HEAD
	buffer, err := ioutil.ReadFile(GIT_DIR + "/HEAD")

	if err != nil {
		return nil, err
	}

	// Convert to string and trim newline
	path := string(buffer[:])
	path = path[:len(path)-1]

	// Get path from HEAD ref
	ref := strings.Split(path, " ")[1]

	// Read commit hash
	buffer, err = ioutil.ReadFile(GIT_DIR + "/" + ref)

	if err != nil {
		return nil, err
	}

	// Convert to string
	hash := string(buffer[:])

	// Make short hash
	hash = hash[:7]

	return &hash, nil
}
