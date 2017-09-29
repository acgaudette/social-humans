package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	restart := make(chan bool, 1)
	var server *http.Server

	for {
		select {
		case <-interrupt:
			if server != nil {
				return shutdown(server)
			}

			return nil

		case <-restart:
			go listen(server, restart)

		default:
			if server == nil {
				server = newServer()
				restart <- true
			}
		}
	}
}
