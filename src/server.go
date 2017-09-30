package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func newServer() *http.Server {
	router := NewRouter()
	router.Handle(http.MethodGet, "/", index)

	router.Handle(http.MethodGet, "/login", getLogin)
	router.Handle(http.MethodPost, "/login", login)

	router.Handle(http.MethodGet, "/logout", getLogout)
	router.Handle(http.MethodPost, "/logout", logout)

	router.Handle(http.MethodGet, "/pool", getPool)
	router.Handle(http.MethodPost, "/pool", managePool)

	return &http.Server{
		Addr:    ADDRESS + ":" + PORT,
		Handler: router,
	}
}

func listen(server *http.Server, failure chan bool) {
	log.Printf("Listening on http://%s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Printf("%s", err)
		failure <- true
	}
}

func shutdown(server *http.Server) error {
	log.Printf("Shutting down...")

	background, _ := context.WithTimeout(
		context.Background(), EXIT_TIMEOUT*time.Second,
	)

	return server.Shutdown(background)
}
