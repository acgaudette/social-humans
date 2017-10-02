package handlers

import (
	"../app"
	"../control"
	"../data"
	"../front"
	"net/http"
)

func GetCreate(out http.ResponseWriter, in *http.Request) *app.Error {
	views := control.GetUserAndMakeViews(nil, in)
	return front.ServeTemplate(out, "create", views)
}

func Create(out http.ResponseWriter, in *http.Request) *app.Error {
	// Serve back the page with a status message
	serveStatus := func(status string) *app.Error {
		view := front.StatusView{Status: status}
		views := control.GetUserAndMakeViews(view, in)
		return front.ServeTemplate(out, "create", views)
	}

	/* Read fields from form */

	in.ParseForm()

	handle, err := front.SanitizeFormString("handle", &in.Form)

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
		return front.ServerError(err)
	}

	// Create session and redirect back home
	err = data.AddSession(out, account)
	return front.Redirect("/", err, out, in)
}
