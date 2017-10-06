package handlers

import (
	"../app"
	"../control"
	"../data"
	"net/http"
	"strings"
)

// Get the edit form for a user's post
func GetEditPost(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, _ := data.GetUserFromSession(in)

	// Extract the handle and timestamp from the URL
	tokens := strings.Split(in.URL.Path, "/")
	handle, stamp := tokens[2], tokens[4]

	// Check if user exists
	_, err := data.LoadUser(handle)

	if err != nil {
		return app.NotFound(err)
	}

	// Check if post exists
	post, err := data.LoadPost(handle + "/" + stamp)

	if err != nil {
		return app.NotFound(err)
	}

	// Get active user and build views
	container := control.MakeContainer(active)
	container.SetContent(control.MakePostView(post, active))

	return app.ServeTemplate(out, "edit_post", container)
}

// Update a user's post
func EditPost(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, err := data.GetUserFromSession(in)

	// Redirect to login page if not logged in
	if err != nil {
		return app.Redirect("/login", err, out, in)
	}

	// Extract the handle and timestamp from the URL
	tokens := strings.Split(in.URL.Path, "/")
	handle, stamp := tokens[2], tokens[4]

	// Check if user exists
	_, err = data.LoadUser(handle)

	if err != nil {
		return app.NotFound(err)
	}

	// Check if post exists
	post, err := data.LoadPost(handle + "/" + stamp)

	if err != nil {
		return app.NotFound(err)
	}

	// Serve back the page with a status message
	serveStatus := func(message string) *app.Error {
		container := control.MakeContainer(active)
		container.SetContent(control.MakePostView(post, active))
		container.SetStatus(control.MakeStatusView(message))

		return app.ServeTemplate(out, "edit_post", container)
	}

	/* Read fields from form */

	in.ParseForm()

	title, content, status := control.ReadPostForm(&in.Form)

	if status != nil {
		return serveStatus(*status)
	}

	// Update post and redirect
	if err = data.UpdatePost(*title, *content, handle, stamp); err != nil {
		return app.ServerError(err)
	}

	return app.Redirect("/", nil, out, in)
}
