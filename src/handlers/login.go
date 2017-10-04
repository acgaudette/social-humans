package handlers

import (
	"../app"
	"../control"
	"../data"
	"log"
	"net/http"
)

func GetLogin(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, _ := data.GetUserFromSession(in)

	// Serve template with no content
	views := control.MakeContainer(nil, nil, active)
	return app.ServeTemplate(out, "login", views)
}

func Login(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, _ := data.GetUserFromSession(in)

	// Serve back the page with a status message
	serveStatus := func(message string) *app.Error {
		status := control.MakeStatusView(message)
		views := control.MakeContainer(nil, status, active)
		return app.ServeTemplate(out, "login", views)
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

	password := in.Form.Get("password")

	if password == "" {
		return serveStatus("Password required!")
	}

	// Load user account
	account, err := data.LoadUser(handle)

	if err != nil {
		log.Printf("%s", err)
		return serveStatus("User does not exist!")

		// Validate password
	} else if err = account.Validate(password); err != nil {
		log.Printf("%s", err)
		return serveStatus("Invalid password")
	}

	// Join existing user session (if it exists) and redirect back home
	err = data.JoinSession(out, account)
	return app.Redirect("/", err, out, in)
}
