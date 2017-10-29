package handlers

import (
	"../../smhb"
	"../app"
	"../control"
	"../data"
	"log"
	"net/http"
)

func GetPool(out http.ResponseWriter, in *http.Request) *app.Error {
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

	// Initialize an empty status view
	status := control.MakeStatusView("")

	// Get the pool view for the current user
	view, err := control.MakePoolView(active.Handle(), *token)

	if err != nil {
		// Update status message with regards to the error
		switch err.(type) {
		case *control.EmptyPoolError:
			status.Status = "Your pool is empty!"

		default:
			log.Printf("%s", err)
		}
	}

	// Build views
	container := control.MakeContainer()
	container.SetActive(control.MakeActiveView(active))
	container.SetContent(view)
	container.SetStatus(status)

	// Serve
	return app.ServeTemplate(out, "pool", container)
}

func ManagePool(out http.ResponseWriter, in *http.Request) *app.Error {
	active, token, err := data.GetUserFromSession(in)

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
		// Get users slice from pool view
		view, err := control.MakePoolView(active.Handle(), *token)

		if err != nil {
			// Update status message with regards to the error
			switch err.(type) {
			case *control.EmptyPoolError:
				if message == "" {
					message = "Your pool is empty!"
				}

			default:
				log.Printf("%s", err)
			}
		}

		// Build views
		container := control.MakeContainer()
		container.SetActive(control.MakeActiveView(active))
		container.SetContent(view)
		container.SetStatus(control.MakeStatusView(message))

		// Serve
		return app.ServeTemplate(out, "pool", container)
	}

	/* Read fields from form */

	in.ParseForm()

	// Target user handle to operate on
	target, err := control.SanitizeFormString("handle", &in.Form)

	if err != nil {
		return serveStatus("Invalid username")
	}

	if target == "" {
		return serveStatus("Target username required!")
	}

	// Operation
	action, err := control.ReadFormRadio(
		"action", []string{"add", "block"}, &in.Form,
	)

	if err != nil {
		return serveStatus("Action required!")
	}

	// Load pool from current user
	pool, err := data.Backend.GetPool(active.Handle(), *token)

	if err != nil {
		switch err.(type) {
		case smhb.NotFoundError:
			log.Printf(
				"Pool not found for user \"%s\"!", active.Handle(),
			)

			return serveStatus("Missing pool")
		default:
			return app.ServerError(err)
		}
	}

	// Update user pool
	switch action {
	case "add":
		err = data.Backend.EditPoolAdd(pool.Handle(), target, *token)

	case "block":
		err = data.Backend.EditPoolBlock(pool.Handle(), target, *token)
	}

	if err != nil {
		log.Printf("%s", err)
		return serveStatus("User does not exist!")
	}

	return serveStatus("")
}
