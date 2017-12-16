package data

import (
	"../../smhb"
	"log"
)

// Global backend client
var Backend smhb.Client

// Global backend address
var BackendAddress string
var BackendPort int

// Initialize backend
func init() {
	client, err := smhb.NewClient(
		0,
		smhb.TCP,
	)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	Backend = client
}
