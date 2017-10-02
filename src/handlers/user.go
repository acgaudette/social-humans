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
	// Extract the handle and attempt to load user
	handle := in.URL.Path[strings.LastIndex(in.URL.Path, "/")+1:]
	account, err := data.LoadUser(handle)

	// User does not exist
	if err != nil {
		return front.NotFound(err)
	}

	view := control.GetUserView(account, "", in)
	views := control.GetUserAndMakeViews(view, in)
	return front.ServeTemplate(out, "user", views)
}
