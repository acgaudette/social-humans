package handlers

import (
	"../app"
	"../control"
	"../data"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"
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
	post, err := data.LoadPost(handle + "/" + stamp)

	if err != nil {
		return app.NotFound(err)
	}

	// Load current user, if available
	active, _ := data.GetUserFromSession(in)

	// Build views
	container := control.MakeContainer(active)
	container.SetContent(control.MakePostView(post, active))

	// Serve
	return app.ServeTemplate(out, "post", container)
}
