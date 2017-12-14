package main

import (
	"flag"
)

var Port int

// Parse command-line options
func init() {
	flag.IntVar(
		&Port,
		"p",
		1234,
		"The port the server will run on",
	)

	flag.Parse()
}
