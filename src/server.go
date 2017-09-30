package main

import (
	"./handlers"
	"context"
	"log"
	"net/http"
	"time"
)

func newServer() *http.Server {
	mux := newRouter()
	mux.handle(http.MethodGet, "/", handlers.Index)

	mux.handle(http.MethodGet, "/login", handlers.GetLogin)
	mux.handle(http.MethodPost, "/login", handlers.Login)

	mux.handle(http.MethodGet, "/logout", handlers.GetLogout)
	mux.handle(http.MethodPost, "/logout", handlers.Logout)

	mux.handle(http.MethodGet, "/pool", handlers.GetPool)
	mux.handle(http.MethodPost, "/pool", handlers.ManagePool)

	mux.handle(http.MethodGet, "/user/*", handlers.GetUser)

	return &http.Server{
		Addr:    ADDRESS + ":" + PORT,
		Handler: mux,
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
