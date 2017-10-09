package data

import (
	"flag"
	"log"
)

var CommitHash *string

// Parse command-line options
func init() {
	var gitDir = flag.String(
		"git-root",
		DEFAULT_GIT_ROOT,
		"base git directory for the server",
	)

	flag.Parse()

	/* Load commit hash */

	hash, err := getCommitHash(*gitDir)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	CommitHash = hash
}
