package handlers

import (
	"../app"
	"../control"
	"../data"
	"../front"
	"net/http"
)

func Index(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	account, err := data.GetUserFromSession(in)

	// User exists, so serve with user view
	if err == nil {
		view := control.GetUserView(account, "", in)
		return front.ServeTemplate(out, "index", view)
	}

	// Otherwise, serve plain template
	return front.ServeTemplate(out, "index", nil)
}
