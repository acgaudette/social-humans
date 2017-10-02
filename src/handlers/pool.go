package handlers

import (
	"../app"
	"../control"
	"../data"
	"../front"
	"log"
	"net/http"
)

func GetPool(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	account, err := data.GetUserFromSession(in)

	// Redirect to login page if there is no session open
	if err != nil {
		return front.Redirect("/login", err, out, in)
	}

	// Get the pool view for the current user
	view, err := control.GetPoolView(account.Handle, "")

	if err != nil {
		log.Printf("%s", err)
	}

	views := control.MakeViews(view, account)
	return front.ServeTemplate(out, "pool", views)
}

func ManagePool(out http.ResponseWriter, in *http.Request) *app.Error {
	account, err := data.GetUserFromSession(in)

	if err != nil {
		return front.Redirect("/login", err, out, in)
	}

	// Serve back the page with a status message
	serveStatus := func(status string) *app.Error {
		// Get users slice from pool view
		view, err := control.GetPoolView(account.Handle, status)

		if err != nil {
			log.Printf("%s", err)
		}

		views := control.MakeViews(view, account)
		return front.ServeTemplate(out, "pool", views)
	}

	/* Read fields from form */

	in.ParseForm()

	// Target user handle to operate on
	target, err := front.SanitizeFormString("handle", &in.Form)

	if err != nil {
		return serveStatus("Invalid username")
	}

	if target == "" {
		return serveStatus("Target username required!")
	}

	// Operation
	action, err := front.ReadFormRadio(
		"action", []string{"add", "block"}, &in.Form,
	)

	if err != nil {
		return serveStatus("Action required!")
	}

	// Load pool from current user
	pool, err := data.LoadPool(account.Handle)

	if err != nil {
		log.Printf("Pool not found for user \"%s\"! Rebuilding...", account.Handle)

		pool, err = data.AddPool(account.Handle)

		if err != nil {
			return front.ServerError(err)
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
