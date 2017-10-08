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
	active, _ := data.GetUserFromSession(in)

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
