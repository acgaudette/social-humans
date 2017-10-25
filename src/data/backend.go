package data

import (
	"../../smhb"
)

// Global backend client
var backend smhb.Client

// Initialize backend
func init() {
	backend = smhb.NewClient(
		BACKEND_ADDRESS,
		BACKEND_PORT,
		smhb.TCP,
	)
}
