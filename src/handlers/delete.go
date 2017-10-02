package handlers

import (
	"../app"
	"../data"
	"../front"
	"net/http"
)

func GetDelete(out http.ResponseWriter, in *http.Request) *app.Error {
	// Redirect to user page
	return Me(out, in)
}

func Delete(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, err := data.GetUserFromSession(in)

	// Redirect to login page if there is no session open
	if err != nil {
		front.Redirect("/login", err, out, in)
	}

	// Remove user and logout
	data.RemoveUser(active.Handle)
	return Logout(out, in)
}
