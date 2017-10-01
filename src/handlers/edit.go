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

	view := control.GetUserView(account, "", in)
	return front.ServeTemplate(out, "edit", view)
}

func Edit(out http.ResponseWriter, in *http.Request) *app.Error {
	account, err := data.GetUserFromSession(in)

	if err != nil {
		front.Redirect("/login", err, out, in)
	}

	// Serve back the page with a status message
	serveStatus := func(status string) *app.Error {
		view := control.GetUserView(account, status, in)
		return front.ServeTemplate(out, "edit", view)
	}

	/* Read fields from form */

	in.ParseForm()

	name, err := front.ReadFormString("name", false, &in.Form)

	if err != nil {
		log.Printf("%s", err)
	}

	if name != "" {
		// Set new full name for user
		if err = account.SetName(name); err != nil {
			return &app.Error{
				Native: err,
				Code:   app.SERVER,
			}
		}
	}

	old, err := front.ReadFormString("oldPassword", false, &in.Form)

	if err != nil {
		log.Printf("%s", err)
	}

	password, err := front.ReadFormString("newPassword", false, &in.Form)

	if err != nil {
		log.Printf("%s", err)
	}

	confirm, err := front.ReadFormString("confirmPassword", false, &in.Form)

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

	// No errors
	return serveStatus("Information updated")
}
