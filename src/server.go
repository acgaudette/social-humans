package main

import (
	"./handlers"
	"context"
	"log"
	"net/http"
	"time"
)

func newServer() *http.Server {
	router := NewRouter()
	router.Handle(http.MethodGet, "/", handlers.Index)

	router.Handle(http.MethodGet, "/login", handlers.GetLogin)
	router.Handle(http.MethodPost, "/login", handlers.Login)

	router.Handle(http.MethodGet, "/logout", handlers.GetLogout)
	router.Handle(http.MethodPost, "/logout", handlers.Logout)

	router.Handle(http.MethodGet, "/pool", handlers.GetPool)
	router.Handle(http.MethodPost, "/pool", handlers.ManagePool)

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
