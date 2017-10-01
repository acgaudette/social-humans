package handlers

import (
	"../app"
	"../data"
	"../front"
	"net/http"
)

func GetLogout(out http.ResponseWriter, in *http.Request) *app.Error {
	return front.Redirect("/", nil, out, in)
}

func Logout(out http.ResponseWriter, in *http.Request) *app.Error {
	data.ClearSession(out)
	return front.Redirect("/", nil, out, in)
}
