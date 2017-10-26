package data

import (
	"../../smhb"
)

// Global backend client
var Backend smhb.Client

// Initialize backend
func init() {
	Backend = smhb.NewClient(
		BACKEND_ADDRESS,
		BACKEND_PORT,
		smhb.TCP,
	)
}
