package handlers

import (
	"../app"
	"../control"
	"../data"
	"../front"
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
	views := control.MakeViews(view, account)
	return front.ServeTemplate(out, "edit", views)
}

func Edit(out http.ResponseWriter, in *http.Request) *app.Error {
	account, err := data.GetUserFromSession(in)

	if err != nil {
		front.Redirect("/login", err, out, in)
	}

	// Serve back the page with a status message
	serveStatus := func(status string) *app.Error {
		view := control.GetUserView(account, status, in)
		views := control.MakeViews(view, account)
		return front.ServeTemplate(out, "edit", views)
	}

	/* Read fields from form */

	in.ParseForm()

	name := in.Form.Get("name")
	if name != "" {
		// Set new full name for user
		if err = account.SetName(name); err != nil {
			return front.ServerError(err)
		}
	}

	old := in.Form.Get("oldPassword")
	password := in.Form.Get("newPassword")
	confirm := in.Form.Get("confirmPassword")

	// Check if new passwords match and are valid
	if password == confirm && password != "" {
		// Validate password
		if err = account.Validate(old); err != nil {
			return serveStatus("Incorrect password")
		}

		// Set new user password
		if err = account.UpdatePassword(password); err != nil {
			return front.ServerError(err)
		}

		// New password doesn't match the confirmation
	} else if password != confirm {
		return serveStatus("Passwords don't match!")

		// No input to the form
	} else if name == "" && old == "" && password == "" && confirm == "" {
		return serveStatus("No input")

		// Old password was not given
	} else if old != "" {
		return serveStatus("No password supplied")
	}

	// No errors
	return serveStatus("Information updated")
}
