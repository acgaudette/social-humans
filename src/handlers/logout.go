package handlers

import (
	"../app"
	"../data"
	"net/http"
)

func GetLogout(out http.ResponseWriter, in *http.Request) *app.Error {
	// Redirect logout route to root
	return app.Redirect("/", nil, out, in)
}

func Logout(out http.ResponseWriter, in *http.Request) *app.Error {
	// Clear session and log out
	data.ClearSession(out)
	return app.Redirect("/", nil, out, in)
}
