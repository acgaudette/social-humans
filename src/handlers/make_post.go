package handlers

import (
	"../app"
	"../control"
	"../data"
	"fmt"
	"net/http"
	"unicode/utf8"
)

func GetMakePost(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, err := data.GetUserFromSession(in)

	// Redirect to login page if there is no session open
	if err != nil {
		return app.Redirect("/login", err, out, in)
	}

	container := control.MakeContainer(active)
	return app.ServeTemplate(out, "make_post", container)
}

func MakePost(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, err := data.GetUserFromSession(in)

	// Redirect to login page if not logged in
	if err != nil {
		return app.Redirect("/login", err, out, in)
	}

	// Serve back the page with a status message
	serveStatus := func(message string) *app.Error {
		container := control.MakeContainer(active)
		container.SetStatus(control.MakeStatusView(message))
		return app.ServeTemplate(out, "make_post", container)
	}

	/* Read fields from form */

	in.ParseForm()

	title := in.Form.Get("title")

	if title == "" {
		return serveStatus("Title required!")
	}

	// Check title character limit
	if utf8.RuneCountInString(title) > data.TITLE_LIMIT {
		return serveStatus(
			fmt.Sprintf(
				"Post title must be under %v characters", data.TITLE_LIMIT,
			),
		)
	}

	content := in.Form.Get("content")

	if content == "" {
		return serveStatus("Post content required!")
	}

	// Check content character limit
	if utf8.RuneCountInString(content) > data.CONTENT_LIMIT {
		return serveStatus(
			fmt.Sprintf(
				"Post content must be under %v characters", data.CONTENT_LIMIT,
			),
		)
	}

	// Create new post
	err = data.NewPost(title, content, active.Handle)

	if err != nil {
		return app.ServerError(err)
	}

	// No errors, so go back home
	return app.Redirect("/", nil, out, in)
}
