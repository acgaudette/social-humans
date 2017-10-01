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

	// Redirect to login page if there is no session open
	if err != nil {
		return front.Redirect("/login", err, out, in)
	}

	return front.ServeTemplate(out, "index", control.GetUserView(account, in))
}
