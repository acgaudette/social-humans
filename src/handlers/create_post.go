package handlers

import (
	"../app"
	"../control"
	"../data"
	"net/http"
)

func GetCreatePost(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, err := data.GetUserFromSession(in)

	// Redirect to login page if there is no session open
	if err != nil {
		return app.Redirect("/login", err, out, in)
	}

	container := control.MakeContainer(active)
	return app.ServeTemplate(out, "create_post", container)
}

func CreatePost(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, err := data.GetUserFromSession(in)

	// Redirect to login page if not logged in
	if err != nil {
		return app.Redirect("/login", err, out, in)
	}

	// Serve back the page with a status message
	serveStatus := func(message string) *app.Error {
		container := control.MakeContainer(active)
		container.SetStatus(control.MakeStatusView(message))
		return app.ServeTemplate(out, "create_post", container)
	}

	/* Read fields from form */

	in.ParseForm()

	title, content, appErr := control.ReadPostForm(serveStatus, &in.Form)

	if appErr != nil {
		return appErr
	}

	// Create new post
	err = data.NewPost(title, content, active.Handle)

	if err != nil {
		return app.ServerError(err)
	}

	// No errors, so go back home
	return app.Redirect("/", nil, out, in)
}
