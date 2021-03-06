package handlers

import (
	"../../smhb"
	"../app"
	"../control"
	"../data"
	"fmt"
	"net/http"
	"strings"
)

// Get the edit form for a user's post
func GetEditPost(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, token, err := data.GetUserFromSession(in)

	if err != nil {
		// Connection error
		if _, ok := err.(smhb.ConnectionError); ok {
			return app.ServerError(err)
		}

		// Redirect to login page if there is no session open
		return app.Redirect("/login", err, out, in)
	}

	// Extract the handle and timestamp from the URL
	tokens := strings.Split(in.URL.Path, "/")
	handle, stamp := tokens[2], tokens[4]

	// Check if user exists
	err = data.Backend.CheckUser(handle)

	if err != nil {
		switch err.(type) {
		case smhb.NotFoundError:
			return app.NotFound(err)
		default:
			return app.ServerError(err)
		}
	}

	// Check if post exists
	post, err := data.Backend.GetPost(handle+"/"+stamp, *token)

	if err != nil {
		switch err.(type) {
		case smhb.NotFoundError:
			return app.NotFound(err)
		default:
			return app.ServerError(err)
		}
	}

	// Check active user against post owner
	if handle != active.Handle() {
		return app.Forbidden(
			fmt.Errorf(
				"user \"%s\" attempted to edit post by user \"%s\"",
				active.Handle(), handle,
			),
		)
	}

	// Build views with active user
	container := control.MakeContainer()
	container.SetActive(control.MakeActiveView(active))
	view := control.MakePostView(post, post.WasAuthoredBy(active.Handle()))
	container.SetContent(view)

	// Serve
	return app.ServeTemplate(out, "edit_post", container)
}

// Update a user's post
func EditPost(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, token, err := data.GetUserFromSession(in)

	if err != nil {
		// Connection error
		if _, ok := err.(smhb.ConnectionError); ok {
			return app.ServerError(err)
		}

		// Redirect to login page if there is no session open
		return app.Redirect("/login", err, out, in)
	}

	// Extract the handle and timestamp from the URL
	tokens := strings.Split(in.URL.Path, "/")
	handle, stamp := tokens[2], tokens[4]

	// Check if user exists
	_, err = data.Backend.GetUser(handle)

	if err != nil {
		switch err.(type) {
		case smhb.NotFoundError:
			return app.NotFound(err)
		default:
			return app.ServerError(err)
		}
	}

	// Check if post exists
	post, err := data.Backend.GetPost(handle+"/"+stamp, *token)

	if err != nil {
		switch err.(type) {
		case smhb.NotFoundError:
			return app.NotFound(err)
		default:
			return app.ServerError(err)
		}
	}

	// Check active user against post owner
	if handle != active.Handle() {
		return app.Forbidden(
			fmt.Errorf(
				"user \"%s\" attempted to edit post by user \"%s\"",
				active.Handle(), handle,
			),
		)
	}

	// Serve back the page with a status message
	serveStatus := func(message string) *app.Error {
		container := control.MakeContainer()
		container.SetActive(control.MakeActiveView(active))
		view := control.MakePostView(post, post.WasAuthoredBy(active.Handle()))
		container.SetContent(view)
		container.SetStatus(control.MakeStatusView(message))

		return app.ServeTemplate(out, "edit_post", container)
	}

	/* Read fields from form */

	in.ParseForm()

	title, content, status := control.ReadPostForm(&in.Form)

	if status != nil {
		return serveStatus(*status)
	}

	// Update post and redirect
	err = data.Backend.EditPost(post.GetAddress(), *title, *content, *token)

	if err != nil {
		return app.ServerError(err)
	}

	// Redirect back to post
	path := "/user/" + handle + "/post/" + stamp
	return app.Redirect(path, nil, out, in)
}
