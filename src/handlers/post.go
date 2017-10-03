package handlers

import (
	"../app"
	"../control"
	"../data"
	"../front"
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
		return front.NotFound(err)
	}

	// Check if post exists
	post, err := data.LoadPost(handle + "/" + stamp)

	if err != nil {
		return front.NotFound(err)
	}

	// Get active user and build views
	view, active := control.GetUserAndMakePostView(post, in)
	views := control.MakeViews(view, active)

	return front.ServeTemplate(out, "post", views)
}
