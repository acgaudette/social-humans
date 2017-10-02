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

	// Get active user and build views
	view, active := control.GetUserAndMakeUserView(account, "", in)
	views := control.MakeViews(view, active)

	return front.ServeTemplate(out, "user", views)
}
