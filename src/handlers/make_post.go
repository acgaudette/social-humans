package handlers

import (
	"../app"
	"../control"
	"../data"
	"../front"
	"net/http"
	"unicode/utf8"
)

func GetMakePost(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, err := data.GetUserFromSession(in)

	// Redirect to login page if there is no session open
	if err != nil {
		return front.Redirect("/login", err, out, in)
	}

	views := control.MakeViews(nil, active)
	return front.ServeTemplate(out, "make_post", views)
}

func MakePost(out http.ResponseWriter, in *http.Request) *app.Error {
	active, err := data.GetUserFromSession(in)

	if err != nil {
		return front.Redirect("/login", err, out, in)
	}

	// Serve back the page with a status message
	serveStatus := func(status string) *app.Error {
		view := front.StatusView{Status: status}
		views := control.MakeViews(view, active)
		return front.ServeTemplate(out, "make_post", views)
	}

	/* Read fields from form */

	in.ParseForm()

	title := in.Form.Get("title")

	if title == "" {
		return serveStatus("Title required!")
	}

	if utf8.RuneCountInString(title) > 20 {
		return serveStatus("Post title must be under 20 characters")
	}

	content := in.Form.Get("content")

	if content == "" {
		return serveStatus("Post content required!")
	}

	if utf8.RuneCountInString(content) > 100 {
		return serveStatus("Post content must be under 100 characters")
	}

	err = data.NewPost(title, content, active.Handle)

	if err != nil {
		return front.ServerError(err)
	}

	// No errors, so go back home
	return front.Redirect("/", nil, out, in)
}
