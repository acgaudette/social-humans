package main

import "./handlers"

func addRoutes(mux *Router) {
	/* Login page */

	mux.GET("/login", handlers.GetLogin)
	mux.POST("/login", handlers.Login)

	/* Logout route */

	mux.POST("/logout", handlers.Logout)
	mux.GET("/logout", handlers.GetLogout)

	/* Create account page */

	mux.GET("/create", handlers.GetCreate)
	mux.POST("/create", handlers.Create)

	/* Edit account page */

	mux.GET("/edit", handlers.GetEdit)
	mux.POST("/edit", handlers.Edit)

	/* Delete account route */

	mux.POST("/delete", handlers.Delete)
	mux.GET("/delete", handlers.GetDelete)

	/* User pool page */

	mux.GET("/pool", handlers.GetPool)
	mux.POST("/pool", handlers.ManagePool)

	/* User page */

	mux.GET("/user/*", handlers.GetUser)
	mux.GET("/me", handlers.Me) // Redirect

	/* Post page */

	mux.GET("/user/*/post/*", handlers.GetPost)

	/* Post edit page */

	mux.GET("/user/*/post/*/edit", handlers.GetEditPost)
	mux.POST("/user/*/post/*/edit", handlers.EditPost)

	/* Post delete route */

	mux.POST("/user/*/post/*/delete", handlers.DeletePost)

	/* Create post */

	mux.GET("/post", handlers.GetCreatePost)
	mux.POST("/post", handlers.CreatePost)

	/* Stylesheet */

	mux.GET("/style.css", handlers.GetStyle)
}
