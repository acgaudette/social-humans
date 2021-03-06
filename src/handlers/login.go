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
	active, _, err := data.GetUserFromSession(in)

	// Connection error
	if err != nil {
		if _, ok := err.(smhb.ConnectionError); ok {
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
	active, _, err := data.GetUserFromSession(in)

	// Connection error
	if err != nil {
		if _, ok := err.(smhb.ConnectionError); ok {
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

	// Check for existing user
	err = data.Backend.CheckUser(handle)

	if err != nil {
		log.Printf("%s", err)
		switch err.(type) {
		case smhb.NotFoundError:
			return serveStatus("User does not exist!")
		default:
			return serveStatus("Error communicating with server")
		}
	}

	// Join existing user session and redirect back home
	err = data.JoinSession(out, handle, password)

	// Check for errors in authentication
	if err != nil {
		if _, ok := err.(smhb.AuthError); ok {
			return serveStatus("Invalid password")
		}
	}

	return app.Redirect("/", err, out, in)
}
