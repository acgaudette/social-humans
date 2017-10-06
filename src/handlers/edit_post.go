package handlers

import (
	"../app"
	"../control"
	"../data"
	"net/http"
	"strings"
)

// Get the edit form for a user's post
func EditPost(out http.ResponseWriter, in *http.Request) *app.Error {
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
