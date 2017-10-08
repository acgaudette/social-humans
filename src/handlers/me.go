package handlers

import (
	"../app"
	"../data"
	"net/http"
)

func Me(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, err := data.GetUserFromSession(in)

	// Redirect to login page if there is no session open
	if err != nil {
		return app.Redirect("/login", err, out, in)
	}

	// Redirect to user page
	return app.Redirect("/user/"+active.Handle(), nil, out, in)
}
