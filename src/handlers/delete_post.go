package handlers

import (
	"../app"
	"../data"
	"net/http"
	"strings"
)

func DeletePost(out http.ResponseWriter, in *http.Request) *app.Error {
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

	// Delete post and redirect
	if err = data.RemovePost(handle + "/" + stamp); err != nil {
		return app.ServerError(err)
	}

	return app.Redirect("/", nil, out, in)
}