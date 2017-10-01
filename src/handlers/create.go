package handlers

import (
	"../app"
	"../data"
	"../front"
	"log"
	"net/http"
)

func GetCreate(out http.ResponseWriter, in *http.Request) *app.Error {
	return front.ServeTemplate(out, "create", &front.StatusMessage{})
}

func Create(out http.ResponseWriter, in *http.Request) *app.Error {
	in.ParseForm()

	// Serve back the page with a status message
	serveStatus := func(status string) *app.Error {
		message := &front.StatusMessage{Status: status}
		return front.ServeTemplate(out, "create", message)
	}

	handle, err := front.ReadFormString("handle", true, &in.Form)

	if err != nil {
		log.Printf("%s", err)
		return serveStatus("Username required!")
	}

	name, err := front.ReadFormString("name", false, &in.Form)

	if err != nil {
		log.Printf("%s", err)
		return serveStatus("Name required!")
	}

	password, err := front.ReadFormString("password", false, &in.Form)

	if err != nil {
		log.Printf("%s", err)
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
		return &app.Error{
			Native: err,
			Code:   app.SERVER,
		}
	}

	// Create session and redirect back home
	data.AddSession(out, account)
	return front.Redirect("/", nil, out, in)
}
