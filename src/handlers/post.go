package handlers

import (
	"../app"
	"../control"
	"../data"
	"net/http"
	"strings"
)

func GetPost(out http.ResponseWriter, in *http.Request) *app.Error {
	// Extract the handle and timestamp from the URL
	tokens := strings.Split(in.URL.Path, "/")
	handle, stamp := tokens[2], tokens[4]

	// Check if user exists
	_, err := data.LoadUser(handle)

	if err != nil {
		return app.NotFound(err)
	}

	// Check if post exists
	post, err := data.LoadPost(data.BuildPostAddress(handle, stamp))

	if err != nil {
		return app.NotFound(err)
	}

	// Load current user, if available
	active, _ := data.GetUserFromSession(in)

	// Initialize view container
	container := control.MakeContainer()

	// Set relevant fields if a session is active
	if active != nil {
		container.SetActive(control.MakeActiveView(active))
		view := control.MakePostView(post, post.WasAuthoredBy(active.Handle()))
		container.SetContent(view)
	} else {
		container.SetContent(control.MakePostView(post, false))
	}

	// Serve
	return app.ServeTemplate(out, "post", container)
}
