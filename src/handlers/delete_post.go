package handlers

import (
	"../app"
	"../data"
	"fmt"
	"net/http"
	"strings"
)

func GetDeletePost(out http.ResponseWriter, in *http.Request) *app.Error {
	// Redirect back to post
	path := strings.TrimSuffix(in.URL.Path, "/delete")
	return app.Redirect(path, nil, out, in)
}

func DeletePost(out http.ResponseWriter, in *http.Request) *app.Error {
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
	_, err = data.LoadPost(handle + "/" + stamp)

	if err != nil {
		return app.NotFound(err)
	}

	// Check active user against post owner
	if handle != active.Handle {
		return app.Forbidden(
			fmt.Errorf(
				"user \"%s\" attempted to delete post by user \"%s\"",
				active.Handle, handle,
			),
		)
	}

	// Delete post and redirect
	if err = data.RemovePost(handle + "/" + stamp); err != nil {
		return app.ServerError(err)
	}

	return app.Redirect("/", nil, out, in)
}
