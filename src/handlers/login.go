package handlers

import (
	"../app"
	"../data"
	"../front"
	"log"
	"net/http"
)

func GetLogin(out http.ResponseWriter, in *http.Request) *app.Error {
	return front.ServeTemplate(out, "login", &front.StatusMessage{})
}

func Login(out http.ResponseWriter, in *http.Request) *app.Error {
	in.ParseForm()

	serveStatus := func(status string) *app.Error {
		message := &front.StatusMessage{Status: status}
		return front.ServeTemplate(out, "login", message)
	}

	handle, err := front.ReadFormString("handle", true, &in.Form)

	if err != nil {
		log.Printf("%s", err)
		return serveStatus("Username required!")
	}

	password, err := front.ReadFormString("password", false, &in.Form)

	if err != nil {
		log.Printf("%s", err)
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
	data.AddSession(out, account)
	return front.Redirect("/", nil, out, in)
}
