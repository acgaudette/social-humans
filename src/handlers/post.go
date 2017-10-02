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

	in.ParseForm()

	title, err := front.ReadFormString("title", true, &in.Form)
	if err != nil {
		log.Printf("%s", err)
		view.Status = "Title required for post"
		return front.ServeTemplate(out, "post", view)
	}

	content, err := front.ReadFormString("content", true, &in.Form)
	if err != nil {
		log.Printf("%s", err)
		view.Status = "Content required for post"
		return front.ServeTemplate(out, "post", view)
	}

	error := data.SavePost(content, title, view.Handle)
	if error != nil {
		return &app.Error{
			Native: error,
			Code:   app.SERVER,
		}
	} else {
		view.Status = "Successfully created new post"
	}

	return front.ServeTemplate(out, "post", view)
}
