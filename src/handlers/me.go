package handlers

import (
	"../../smhb"
	"../app"
	"../data"
	"net/http"
)

func Me(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, _, err := data.GetUserFromSession(in)

	if err != nil {
		// Connection error
		if _, ok := err.(smhb.ConnectionError); ok {
			return app.ServerError(err)
		}

		// Redirect to login page if there is no session open
		return app.Redirect("/login", err, out, in)
	}

	// Redirect to user page
	return app.Redirect("/user/"+active.Handle(), nil, out, in)
}
