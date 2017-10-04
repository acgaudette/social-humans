package handlers

import (
	"../app"
	"../data"
	"net/http"
)

func GetLogout(out http.ResponseWriter, in *http.Request) *app.Error {
	return app.Redirect("/", nil, out, in)
}

func Logout(out http.ResponseWriter, in *http.Request) *app.Error {
	data.ClearSession(out)
	return app.Redirect("/", nil, out, in)
}
