package handlers

import (
	"../../smhb"
	"../app"
	"../control"
	"../data"
	"log"
	"net/http"
)

func GetLogin(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, err := data.GetUserFromSession(in)

	// Connection error
	if err != nil {
		if _, ok := err.(smhb.NotFoundError); !ok {
			return app.ServerError(err)
		}
	}

	// Initialize view container
	container := control.MakeContainer()

	if active != nil {
		container.SetActive(control.MakeActiveView(active))
	}

	// Serve template with no content
	return app.ServeTemplate(out, "login", container)
}

func Login(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, err := data.GetUserFromSession(in)

	// Connection error
	if err != nil {
		if _, ok := err.(smhb.NotFoundError); !ok {
			return app.ServerError(err)
		}
	}

	// Serve back the page with a status message
	serveStatus := func(message string) *app.Error {
		container := control.MakeContainer()

		if active != nil {
			container.SetActive(control.MakeActiveView(active))
		}

		container.SetStatus(control.MakeStatusView(message))
		return app.ServeTemplate(out, "login", container)
	}

	/* Read fields from form */

	in.ParseForm()

	handle, err := control.SanitizeFormString("handle", &in.Form)

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
	account, err := data.Backend.GetUser(handle)

	if err != nil {
		log.Printf("%s", err)
		switch err.(type) {
		case smhb.NotFoundError:
			return serveStatus("User does not exist!")
		default:
			return serveStatus("Error communicating with server")
		}

		// Validate password
	} else if err = account.Validate(password); err != nil {
		log.Printf("%s", err)
		return serveStatus("Invalid password")
	}

	// Join existing user session (if it exists) and redirect back home
	err = data.JoinSession(out, account)
	return app.Redirect("/", err, out, in)
}
