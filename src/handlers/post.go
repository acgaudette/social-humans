package handlers

import (
	"../app"
	"../data"
	"../front"
	"log"
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

	title, err := front.ReadFormString("title", true, &in.Form)
	if err != nil {
		log.Printf("%s", err)
		serveStatus("Title required for post")
	}

	content, err := front.ReadFormString("content", true, &in.Form)
	if err != nil {
		log.Printf("%s", err)
		serveStatus("Content required for post")
	}

	err = data.NewPost(title, content, account.Handle)
	if err != nil {
		return &app.Error{
			Native: err,
			Code:   app.SERVER,
		}
	}

	return serveStatus("Created new post")
}
