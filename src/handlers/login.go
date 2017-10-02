package handlers

import (
	"../app"
	"../control"
	"../data"
	"../front"
	"log"
	"net/http"
)

func GetLogin(out http.ResponseWriter, in *http.Request) *app.Error {
	views := control.GetUserAndMakeViews(nil, in)
	return front.ServeTemplate(out, "login", views)
}

func Login(out http.ResponseWriter, in *http.Request) *app.Error {
	// Serve back the page with a status message
	serveStatus := func(status string) *app.Error {
		view := front.LoginView{Status: status}
		views := control.GetUserAndMakeViews(view, in)
		return front.ServeTemplate(out, "login", views)
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

	// Create user session and redirect back home
	err = data.JoinSession(out, account)
	return front.Redirect("/", err, out, in)
}
