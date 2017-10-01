package handlers

import (
	"../app"
	"../control"
	"../data"
	"../front"
	"log"
	"net/http"
)

func GetEdit(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	account, err := data.GetUserFromSession(in)

	// Redirect to login page if there is no session open
	if err != nil {
		front.Redirect("/login", err, out, in)
	}

	return front.ServeTemplate(out, "edit", control.GetUserView(account, in))
}

func Edit(out http.ResponseWriter, in *http.Request) *app.Error {
	account, err := data.GetUserFromSession(in)

	if err != nil {
		front.Redirect("/login", err, out, in)
	}

	// Serve back the page with a status message
	serveStatus := func(status string) *app.Error {
		message := control.GetUserView(account, in)
		message.Status = status
		return front.ServeTemplate(out, "edit", &message)
	}

	in.ParseForm()

	name, err := front.ReadFormString("name", false, &in.Form)

	if err != nil {
		log.Printf("%s", err)
	}

	// Set new user full name
	if name != "" {
		if err = account.SetName(name); err != nil {
			return &app.Error{
				Native: err,
				Code:   app.SERVER,
			}
		}
	}

	password, err := front.ReadFormString("newPassword", false, &in.Form)

	if err != nil {
		log.Printf("%s", err)
	}

	confirm, err := front.ReadFormString("confirmPassword", false, &in.Form)

	if err != nil {
		log.Printf("%s", err)
	}

	old, err := front.ReadFormString("oldPassword", false, &in.Form)

	if err != nil {
		log.Printf("%s", err)
	}

	// Check if new passwords match and are valid
	if password == confirm && password != "" {
		// Validate password
		if err = account.Validate(old); err != nil {
			return serveStatus("Invalid password")
		}

		// Set new user password
		if err = account.UpdatePassword(password); err != nil {
			return &app.Error{
				Native: err,
				Code:   app.SERVER,
			}
		}

	} else if password != confirm {
		return serveStatus("Passwords don't match!")
	}

	return serveStatus("Information updated")
}
