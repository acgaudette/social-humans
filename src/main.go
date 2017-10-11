package main

import (
	"log"
	"os"
	"os/signal"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)

	restart := make(chan bool, 1)
	restart <- true

	server := newServer()

	for {
		select {
		// Stop server on keyboard interrupt
		case <-interrupt:
			if server != nil {
				return shutdown(server)
			}

			return nil

		case <-restart:
			go listen(server, restart)
		}
	}
}
