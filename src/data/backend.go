package data

import (
	"../../smhb"
)

// Global backend client
var Backend smhb.Client

// Global backend address
var BackendAddress string
var BackendPort int

// Initialize backend
func init() {
	Backend = smhb.NewClient(
		BackendAddress,
		BackendPort,
		smhb.TCP,
	)
}
