package handlers

import (
	"../app"
	"../control"
	"../data"
	"net/http"
)

func GetCreate(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, _ := data.GetUserFromSession(in)

	// Build views and serve
	views := control.MakeViews(nil, nil, active)
	return app.ServeTemplate(out, "create", views)
}

func Create(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, _ := data.GetUserFromSession(in)

	// Serve back the page with a status message
	serveStatus := func(message string) *app.Error {
		status := control.MakeStatusView(message)
		views := control.MakeViews(nil, status, active)
		return app.ServeTemplate(out, "create", views)
	}

	/* Read fields from form */

	in.ParseForm()

	handle, err := control.SanitizeFormString("handle", &in.Form)

	if err != nil {
		return serveStatus("Invalid username")
	}

	if handle == "" {
		return serveStatus("Username required!")
	}

	name := in.Form.Get("name")

	if name == "" {
		return serveStatus("Name required!")
	}

	password := in.Form.Get("password")

	if password == "" {
		return serveStatus("Password required!")
	}

	// Check for existing user
	account, err := data.LoadUser(handle)

	// If user exists, fail
	if err == nil {
		return serveStatus("Username taken!")
	}

	// Add new user
	account, err = data.AddUser(handle, password, name)

	if err != nil {
		return app.ServerError(err)
	}

	// Create session and redirect back home
	err = data.AddSession(out, account)
	return app.Redirect("/", err, out, in)
}
