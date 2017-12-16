package main

import (
	"../smhb"
	"log"
	"os"
	"os/signal"
	"strconv"
)

const (
	ADDRESS   = "0.0.0.0"
	POOL_SIZE = 8
	DATA_PATH = "data"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// Create data server and serve
func run() error {
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)

	restart := make(chan bool, 1)
	restart <- true

	server := smhb.NewServer(
		ADDRESS,
		Port,
		smhb.TCP,
		POOL_SIZE,
		DATA_PATH+"-"+strconv.Itoa(Port),
	)

	for {
		select {
		case <-interrupt:
			os.Exit(0)

		case <-restart:
			err := server.ListenAndServe()
			log.Printf("%s", err)
			restart <- true
		}
	}
}
