package main

import (
	"./handlers"
	"context"
	"log"
	"net/http"
	"time"
)

func newServer() *http.Server {
	mux := NewRouter(handlers.Index)

	mux.GET("/login", handlers.GetLogin)
	mux.POST("/login", handlers.Login)

	mux.GET("/logout", handlers.GetLogout)
	mux.POST("/logout", handlers.Logout)

	mux.GET("/create", handlers.GetCreate)
	mux.POST("/create", handlers.Create)

	mux.GET("/edit", handlers.GetEdit)
	mux.POST("/edit", handlers.Edit)

	mux.GET("/delete", handlers.GetDelete)
	mux.POST("/delete", handlers.Delete)

	mux.GET("/pool", handlers.GetPool)
	mux.POST("/pool", handlers.ManagePool)

	mux.GET("/me", handlers.Me)
	mux.GET("/user/*", handlers.GetUser)

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
