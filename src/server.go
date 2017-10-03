package main

import (
	"./handlers"
	"context"
	"log"
	"net/http"
	"time"
)

func newServer() *http.Server {
	// Create router and load routes
	mux := NewRouter(handlers.Index)
	addRoutes(mux)

	return &http.Server{
		Addr:    ADDRESS + ":" + PORT,
		Handler: mux,
	}
}

func listen(server *http.Server, failure chan bool) {
	log.Printf("Listening on http://%s", server.Addr)

	// Listen; restart on failure
	if err := server.ListenAndServe(); err != nil {
		log.Printf("%s", err)
		failure <- true
	}
}

// Gracefully shut down server with timeout
func shutdown(server *http.Server) error {
	log.Printf("Shutting down...")

	background, _ := context.WithTimeout(
		context.Background(), EXIT_TIMEOUT*time.Second,
	)

	return server.Shutdown(background)
}
