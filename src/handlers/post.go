package handlers

import (
	"../app"
	"../control"
	"../data"
	"../front"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"
)

func GetPost(out http.ResponseWriter, in *http.Request) *app.Error {
	// Extract the handle and timestamp from the URL
	tokens := strings.Split(in.URL.Path, "/")
	handle, stamp := tokens[2], tokens[4]

	// Check if user exists
	_, err := data.LoadUser(handle)

	if err != nil {
		return front.NotFound(err)
	}

	// Check if post exists
	post, err := data.LoadPost(handle + "/" + stamp)

	if err != nil {
		return front.NotFound(err)
	}

	// Get active user and build views
	view, active := control.GetUserAndMakePostView(post, in)
	views := control.MakeViews(view, active)

	return front.ServeTemplate(out, "post", views)
}

// Gets the edit form for a user's post
func EditPost(out http.ResponseWriter, in *http.Request) *app.Error {
	// Extract the handle and timestamp from the URL
	tokens := strings.Split(in.URL.Path, "/")
	handle, stamp := tokens[2], tokens[4]

	// Check if user exists
	_, err := data.LoadUser(handle)

	if err != nil {
		return front.NotFound(err)
	}

	// Check if post exists
	post, err := data.LoadPost(handle + "/" + stamp)

	if err != nil {
		return front.NotFound(err)
	}

	// Get active user and build views
	view, active := control.GetUserAndMakePostView(post, in)
	views := control.MakeViews(view, active)

	return front.ServeTemplate(out, "edit_post", views)
}

// Updates a user's post
func UpdatePost(out http.ResponseWriter, in *http.Request) *app.Error {
	// Extract the handle and timestamp from the URL
	tokens := strings.Split(in.URL.Path, "/")
	handle, stamp := tokens[2], tokens[4]

	// Check if user exists
	_, err := data.LoadUser(handle)

	if err != nil {
		return front.NotFound(err)
	}

	// Check if post exists
	post, err := data.LoadPost(handle + "/" + stamp)

	if err != nil {
		return front.NotFound(err)
	}

	serveError := func(status string) *app.Error {
		view, active := control.GetUserAndMakePostView(post, in)
		statusView := control.MakeStatusView(status)
		views := control.MakeViewsWithStatus(view, active, statusView)
		return front.ServeTemplate(out, "edit_post", views)
	}

	// Get updated post content
	in.ParseForm()

	title := in.Form.Get("title")

	if title == "" {
		return serveError("Title required!")
	}

	if utf8.RuneCountInString(title) > data.TITLE_LIMIT {
		return serveError(
			fmt.Sprintf(
				"Post title must be under %v characters", data.TITLE_LIMIT,
			),
		)
	}

	content := in.Form.Get("content")

	if content == "" {
		return serveError("Post content required!")
	}

	if utf8.RuneCountInString(content) > data.CONTENT_LIMIT {
		return serveError(
			fmt.Sprintf(
				"Post content must be under %v characters", data.CONTENT_LIMIT,
			),
		)
	}

	// Update post and redirect
	if err = data.UpdatePost(title, content, handle, stamp); err != nil {
		return front.ServerError(err)
	}

	return front.Redirect("/", nil, out, in)
}

func DeletePost(out http.ResponseWriter, in *http.Request) *app.Error {
	// Extract the handle and timestamp from the URL
	tokens := strings.Split(in.URL.Path, "/")
	handle, stamp := tokens[2], tokens[4]

	// Check if user exists
	_, err := data.LoadUser(handle)

	if err != nil {
		return front.NotFound(err)
	}

	// Check if post exists
	_, err = data.LoadPost(handle + "/" + stamp)

	if err != nil {
		return front.NotFound(err)
	}

	// Delete post and redirect
	if err = data.RemovePost(handle + "/" + stamp); err != nil {
		return front.ServerError(err)
	}

	return front.Redirect("/", nil, out, in)
}
