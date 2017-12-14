package main

import (
	"../smhb"
	"log"
)

const (
	ADDRESS   = "0.0.0.0"
	POOL_SIZE = 8
	DATA_PATH = "data"
)

func main() {
	if err := run(); err != nil {
		log.Printf("%s", err)
	}
}

// Create data server and serve
func run() error {
	server := smhb.NewServer(
		ADDRESS,
		Port,
		smhb.TCP,
		POOL_SIZE,
		DATA_PATH,
	)

	for {
		err := server.ListenAndServe()
		log.Printf("%s", err)
	}
}
