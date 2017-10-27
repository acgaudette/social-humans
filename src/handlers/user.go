package handlers

import (
	"../../smhb"
	"../app"
	"../control"
	"../data"
	"net/http"
	"strings"
)

func GetUser(out http.ResponseWriter, in *http.Request) *app.Error {
	// Extract the handle from the URL and attempt to load user
	handle := in.URL.Path[strings.LastIndex(in.URL.Path, "/")+1:]
	account, err := data.Backend.GetUser(handle)

	// User does not exist
	if err != nil {
		switch err.(type) {
		case smhb.NotFoundError:
			return app.NotFound(err)
		default:
			return app.ServerError(err)
		}
	}

	// Load current user, if available
	active, err := data.GetUserFromSession(in)

	// Connection error
	if err != nil {
		if _, ok := err.(smhb.ConnectionError); ok {
			return app.ServerError(err)
		}
	}

	// Initialize view container
	container := control.MakeContainer()

	// Set relevant fields if a session is active
	if active != nil {
		container.SetActive(control.MakeActiveView(active))
		view := control.MakeUserView(account, account.Equals(active))
		container.SetContent(view)
	} else {
		container.SetContent(control.MakeUserView(account, false))
	}

	// Serve
	return app.ServeTemplate(out, "user", container)
}
