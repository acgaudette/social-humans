package handlers

import (
	"../app"
	"../control"
	"../data"
	"../front"
	"net/http"
)

func GetPost(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, err := data.GetUserFromSession(in)

	// Redirect to login page if there is no session open
	if err != nil {
		return front.Redirect("/login", err, out, in)
	}

	views := control.MakeViews(nil, active)
	return front.ServeTemplate(out, "post", views)
}

func CreatePost(out http.ResponseWriter, in *http.Request) *app.Error {
	active, err := data.GetUserFromSession(in)

	if err != nil {
		return front.Redirect("/login", err, out, in)
	}

	// Serve back the page with a status message
	serveStatus := func(status string) *app.Error {
		view := front.StatusView{Status: status}
		views := control.MakeViews(view, active)
		return front.ServeTemplate(out, "post", views)
	}

	/* Read fields from form */

	in.ParseForm()

	title := in.Form.Get("title")

	if title == "" {
		return serveStatus("Title required!")
	}

	content := in.Form.Get("content")

	if content == "" {
		return serveStatus("Post content required!")
	}

	err = data.NewPost(title, content, active.Handle)

	if err != nil {
		return front.ServerError(err)
	}

	// No errors, so go back home
	return front.Redirect("/", nil, out, in)
}
