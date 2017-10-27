package handlers

import (
	"../../smhb"
	"../app"
	"../control"
	"../data"
	"net/http"
)

func GetEdit(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, err := data.GetUserFromSession(in)

	if err != nil {
		// Connection error
		if _, ok := err.(smhb.ConnectionError); ok {
			return app.ServerError(err)
		}

		// Redirect to login page if there is no session open
		return app.Redirect("/login", err, out, in)
	}

	// Otherwise, build views with active user
	container := control.MakeContainer()
	container.SetActive(control.MakeActiveView(active))
	container.SetContent(control.MakeUserView(active, true))

	// Serve
	return app.ServeTemplate(out, "edit", container)
}

func Edit(out http.ResponseWriter, in *http.Request) *app.Error {
	active, err := data.GetUserFromSession(in)

	if err != nil {
		// Connection error
		if _, ok := err.(smhb.ConnectionError); ok {
			return app.ServerError(err)
		}

		// Redirect to login page if there is no session open
		return app.Redirect("/login", err, out, in)
	}

	// Serve back the page with a status message
	serveStatus := func(message string) *app.Error {
		container := control.MakeContainer()

		// Active user is guaranteed to not be nil
		container.SetActive(control.MakeActiveView(active))

		container.SetContent(control.MakeUserView(active, true))
		container.SetStatus(control.MakeStatusView(message))

		return app.ServeTemplate(out, "edit", container)
	}

	/* Read fields from form */

	in.ParseForm()

	name := in.Form.Get("name")
	if name != "" {
		// Set new full name for user
		if err = data.Backend.EditUserName(active.Handle(), name); err != nil {
			return app.ServerError(err)
		}

		// Reload updated user
		active, err = data.Backend.GetUser(active.Handle())

		if err != nil {
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
		err = data.Backend.EditUserPassword(active.Handle(), password)

		if err != nil {
			return app.ServerError(err)
		}

		// Reload updated user
		active, err = data.Backend.GetUser(active.Handle())

		if err != nil {
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
