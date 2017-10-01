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
		return &app.Error{
			Native: err,
			Code:   app.NOT_FOUND,
		}
	}

	view := control.GetUserView(account, "", in)
	return front.ServeTemplate(out, "user", view)
}
