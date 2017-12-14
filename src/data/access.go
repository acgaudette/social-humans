package data

import (
	"../../smhb"
)

// Global access client (used for sessions)
var access smhb.Access
var accessContext smhb.ServerContext

// Initialize access
func init() {
	access = smhb.NewFileAccess()
	accessContext = smhb.NewServerContext(DATA_PATH, BackendAddress, BackendPort)
}
