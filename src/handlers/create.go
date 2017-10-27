package handlers

import (
	"../../smhb"
	"../app"
	"../control"
	"../data"
	"net/http"
)

func GetCreate(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, err := data.GetUserFromSession(in)

	// Connection error
	if err != nil {
		if _, ok := err.(smhb.NotFoundError); !ok {
			return app.ServerError(err)
		}
	}

	// Initialize view container
	container := control.MakeContainer()

	if active != nil {
		container.SetActive(control.MakeActiveView(active))
	}

	// Serve template with no content
	return app.ServeTemplate(out, "create", container)
}

func Create(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, err := data.GetUserFromSession(in)

	// Connection error
	if err != nil {
		if _, ok := err.(smhb.NotFoundError); !ok {
			return app.ServerError(err)
		}
	}

	// Serve back the page with a status message
	serveStatus := func(message string) *app.Error {
		container := control.MakeContainer()

		if active != nil {
			container.SetActive(control.MakeActiveView(active))
		}

		container.SetStatus(control.MakeStatusView(message))

		return app.ServeTemplate(out, "create", container)
	}

	/* Read fields from form */

	in.ParseForm()

	handle, name, password, status := control.ReadCreateForm(&in.Form)

	if status != nil {
		return serveStatus(*status)
	}

	// Check for existing user
	account, err := data.Backend.GetUser(*handle)

	// If user exists, fail
	if err == nil {
		return serveStatus("Username taken!")
	} else {
		// Check for connection error
		if _, ok := err.(smhb.NotFoundError); !ok {
			return serveStatus("Error communicating with server")
		}
	}

	// Add new user
	account, err = data.Backend.AddUser(*handle, *password, *name)

	if err != nil {
		return app.ServerError(err)
	}

	// Create session and redirect back home
	err = data.AddSession(out, account)
	return app.Redirect("/", err, out, in)
}
