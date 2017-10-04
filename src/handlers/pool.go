package handlers

import (
	"../app"
	"../control"
	"../data"
	"log"
	"net/http"
)

func GetPool(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, err := data.GetUserFromSession(in)

	// Redirect to login page if there is no session open
	if err != nil {
		return app.Redirect("/login", err, out, in)
	}

	// Initialize an empty status view
	status := control.MakeStatusView("")

	// Get the pool view for the current user
	view, err := control.MakePoolView(active.Handle)

	if err != nil {
		// Update status message with regards to the error
		switch err.(type) {
		case *control.EmptyPoolError:
			status.Status = "Your pool is empty!"

		default:
			log.Printf("%s", err)
		}
	}

	views := control.MakeViews(view, status, active)
	return app.ServeTemplate(out, "pool", views)
}

func ManagePool(out http.ResponseWriter, in *http.Request) *app.Error {
	active, err := data.GetUserFromSession(in)

	if err != nil {
		return app.Redirect("/login", err, out, in)
	}

	// Serve back the page with a status message
	serveStatus := func(message string) *app.Error {
		// Get users slice from pool view
		view, err := control.MakePoolView(active.Handle)

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

		status := control.MakeStatusView(message)
		views := control.MakeViews(view, status, active)
		return app.ServeTemplate(out, "pool", views)
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
	pool, err := data.LoadPool(active.Handle)

	if err != nil {
		log.Printf("Pool not found for user \"%s\"! Rebuilding...", active.Handle)

		pool, err = data.AddPool(active.Handle)

		if err != nil {
			return app.ServerError(err)
		}
	}

	// Update user pool
	switch action {
	case "add":
		err = pool.Add(target)

	case "block":
		err = pool.Block(target)
	}

	if err != nil {
		log.Printf("%s", err)
		return serveStatus("User does not exist!")
	}

	return serveStatus("")
}
