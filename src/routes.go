package main

import "./handlers"

func addRoutes(mux *Router) {
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

	mux.GET("/post", handlers.GetPost)
	mux.POST("/post", handlers.CreatePost)
}
