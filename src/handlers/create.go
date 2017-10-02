package handlers

import (
	"../app"
	"../data"
	"../front"
	"log"
	"net/http"
)

func GetCreate(out http.ResponseWriter, in *http.Request) *app.Error {
	view := &front.LoginView{}

	// Load current user, if available
	account, err := data.GetUserFromSession(in)

	// Fill view
	if err == nil {
		view.Handle = account.Handle
		view.IsLoggedIn = true
	}

	return front.ServeTemplate(out, "create", view)
}

func Create(out http.ResponseWriter, in *http.Request) *app.Error {
	view := &front.LoginView{}
	active, err := data.GetUserFromSession(in)

	if err == nil {
		view.Handle = active.Handle
		view.IsLoggedIn = true
	}

	// Serve back the page with a status message
	serveStatus := func(status string) *app.Error {
		view.Status = status
		return front.ServeTemplate(out, "create", view)
	}

	/* Read fields from form */

	in.ParseForm()

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
		return front.ServerError(err)
	}

	// Create session and redirect back home
	err = data.AddSession(out, account)
	return front.Redirect("/", err, out, in)
}
