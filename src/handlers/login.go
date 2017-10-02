package handlers

import (
	"../app"
	"../data"
	"../front"
	"log"
	"net/http"
)

func GetLogin(out http.ResponseWriter, in *http.Request) *app.Error {
	view := &front.LoginView{}

	// Load current user, if available
	account, err := data.GetUserFromSession(in)

	// Fill view
	if err == nil {
		view.Handle = account.Handle
		view.IsLoggedIn = true
	}

	return front.ServeTemplate(out, "login", view)
}

func Login(out http.ResponseWriter, in *http.Request) *app.Error {
	view := &front.LoginView{}
	account, err := data.GetUserFromSession(in)

	if err == nil {
		view.Handle = account.Handle
		view.IsLoggedIn = true
	}

	/* Read fields from form */

	// Serve back the page with a status message
	serveStatus := func(status string) *app.Error {
		view.Status = status
		return front.ServeTemplate(out, "login", view)
	}

	in.ParseForm()

	handle, err := front.SanitizeFormString("handle", &in.Form)

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
	account, err = data.LoadUser(handle)

	if err != nil {
		log.Printf("%s", err)
		return serveStatus("User does not exist!")

		// Validate password
	} else if err = account.Validate(password); err != nil {
		log.Printf("%s", err)
		return serveStatus("Invalid password")
	}

	// Create user session and redirect back home
	err = data.AddSession(out, account)
	return front.Redirect("/", err, out, in)
}
