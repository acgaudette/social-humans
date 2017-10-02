package handlers

import (
	"../app"
	"../data"
	"../front"
	"net/http"
)

func GetPost(out http.ResponseWriter, in *http.Request) *app.Error {
	view := &front.PostView{}

	// Load current user, if available
	account, err := data.GetUserFromSession(in)

	// Fill view
	if err == nil {
		view.Handle = account.Handle
		view.Status = "" // No error message
	}

	return front.ServeTemplate(out, "post", view)
}

func CreatePost(out http.ResponseWriter, in *http.Request) *app.Error {
	view := &front.PostView{}

	account, err := data.GetUserFromSession(in)
	if err != nil {
		return front.Redirect("/login", err, out, in)
	}
	view.Handle = account.Handle

	// Serve back the page with a status message
	serveStatus := func(status string) *app.Error {
		view.Status = status
		return front.ServeTemplate(out, "post", view)
	}

	/* Read fields from form */

	in.ParseForm()

	title := in.Form.Get("title")
	if title == "" {
		serveStatus("Title content required!")
	}

	content := in.Form.Get("content")
	if content == "" {
		serveStatus("Post content required!")
	}

	err = data.NewPost(title, content, account.Handle)
	if err != nil {
		return front.ServerError(err)
	}

	return serveStatus("Created new post")
}
