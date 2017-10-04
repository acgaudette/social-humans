package handlers

import (
	"../app"
	"../control"
	"../data"
	"../front"
	"net/http"
	"strings"
)

func GetUser(out http.ResponseWriter, in *http.Request) *app.Error {
	// Extract the handle from the URL and attempt to load user
	handle := in.URL.Path[strings.LastIndex(in.URL.Path, "/")+1:]
	account, err := data.LoadUser(handle)

	// User does not exist
	if err != nil {
		return front.NotFound(err)
	}

	// Load current user, if available
	active, _ := data.GetUserFromSession(in)

	// Build views and serve
	view := control.MakeUserView(account, active)
	views := control.MakeViews(view, active)
	return front.ServeTemplate(out, "user", views)
}
