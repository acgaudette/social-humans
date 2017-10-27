package handlers

import (
	"../../smhb"
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
	_, err := data.Backend.GetUser(handle)

	if err != nil {
		switch err.(type) {
		case smhb.NotFoundError:
			return app.NotFound(err)
		default:
			return app.ServerError(err)
		}
	}

	// Check if post exists
	post, err := data.Backend.GetPost(handle + "/" + stamp)

	if err != nil {
		switch err.(type) {
		case smhb.NotFoundError:
			return app.NotFound(err)
		default:
			return app.ServerError(err)
		}
	}

	// Load current user, if available
	active, err := data.GetUserFromSession(in)

	// Connection error
	if err != nil {
		if _, ok := err.(smhb.NotFoundError); !ok {
			return app.ServerError(err)
		}
	}

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
