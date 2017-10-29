package handlers

import (
	"../../smhb"
	"../app"
	"../control"
	"../data"
	"net/http"
)

func GetCreatePost(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, _, err := data.GetUserFromSession(in)

	if err != nil {
		// Connection error
		if _, ok := err.(smhb.ConnectionError); ok {
			return app.ServerError(err)
		}

		// Redirect to login page if there is no session open
		return app.Redirect("/login", err, out, in)
	}

	// Otherwise, create container with active user and serve
	container := control.MakeContainer()
	container.SetActive(control.MakeActiveView(active))
	return app.ServeTemplate(out, "create_post", container)
}

func CreatePost(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, token, err := data.GetUserFromSession(in)

	if err != nil {
		// Connection error
		if _, ok := err.(smhb.ConnectionError); ok {
			return app.ServerError(err)
		}

		// Redirect to login page if there is no session open
		return app.Redirect("/login", err, out, in)
	}

	// Serve back the page with a status message
	serveStatus := func(message string) *app.Error {
		container := control.MakeContainer()

		if active != nil {
			container.SetActive(control.MakeActiveView(active))
		}

		container.SetStatus(control.MakeStatusView(message))
		return app.ServeTemplate(out, "create_post", container)
	}

	/* Read fields from form */

	in.ParseForm()

	title, content, status := control.ReadPostForm(&in.Form)

	if status != nil {
		return serveStatus(*status)
	}

	// Create new post
	err = data.Backend.AddPost(*title, *content, active.Handle(), *token)

	if err != nil {
		return app.ServerError(err)
	}

	// No errors, so go back home
	return app.Redirect("/", nil, out, in)
}
