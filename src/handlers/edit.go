package handlers

import (
	"../app"
	"../control"
	"../data"
	"net/http"
)

func GetEdit(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, err := data.GetUserFromSession(in)

	// Redirect to login page if there is no session open
	if err != nil {
		app.Redirect("/login", err, out, in)
	}

	// Build views
	container := control.MakeContainer(active)
	container.SetContent(control.MakeUserView(active, active))

	// Serve
	return app.ServeTemplate(out, "edit", container)
}

func Edit(out http.ResponseWriter, in *http.Request) *app.Error {
	active, err := data.GetUserFromSession(in)

	if err != nil {
		app.Redirect("/login", err, out, in)
	}

	// Serve back the page with a status message
	serveStatus := func(message string) *app.Error {
		container := control.MakeContainer(active)
		container.SetContent(control.MakeUserView(active, active))
		container.SetStatus(control.MakeStatusView(message))

		return app.ServeTemplate(out, "edit", container)
	}

	/* Read fields from form */

	in.ParseForm()

	name := in.Form.Get("name")
	if name != "" {
		// Set new full name for user
		if err = active.SetName(name); err != nil {
			return app.ServerError(err)
		}
	}

	old := in.Form.Get("oldPassword")
	password := in.Form.Get("newPassword")
	confirm := in.Form.Get("confirmPassword")

	// Check if new passwords match and are valid
	if password == confirm && password != "" {
		// Validate password
		if err = active.Validate(old); err != nil {
			return serveStatus("Incorrect password")
		}

		// Set new user password
		if err = active.UpdatePassword(password); err != nil {
			return app.ServerError(err)
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
