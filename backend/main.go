package main

import (
	"../smhb"
	"log"
)

const (
	ADDRESS   = "0.0.0.0"
	PORT      = 19138
	DATA_PATH = "data"
)

func main() {
	if err := run(); err != nil {
		log.Printf("%s", err)
	}
}

func run() error {
	server := smhb.NewServer(
		ADDRESS,
		PORT,
		smhb.TCP,
		DATA_PATH,
	)

	return server.ListenAndServe()
}
