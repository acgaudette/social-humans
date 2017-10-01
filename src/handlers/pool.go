package handlers

import (
	"../control"
	"../data"
	"../front"
	"log"
	"net/http"
)

func GetPool(writer http.ResponseWriter, request *http.Request) {
	// Load current user, if available
	account, err := data.GetUserFromSession(request)

	// Redirect to login page if there is no session open
	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusFound)
		log.Printf("%s", err)
		return
	}

	// Get the pool view for the current user
	users, err := control.GetPoolView(account.Handle, "")

	if err != nil {
		log.Printf("%s", err)
	}

	err = front.ServeTemplate(writer, "pool", users)

	if err != nil {
		log.Printf("%s", err)
	}
}

func ManagePool(writer http.ResponseWriter, request *http.Request) {
	account, err := data.GetUserFromSession(request)

	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusFound)
		log.Printf("%s", err)
		return
	}

	request.ParseForm()

	serveStatus := func(status string) {
		// Get users slice from pool view
		// Note: users will never be nil
		users, err := control.GetPoolView(account.Handle, status)

		if err != nil {
			log.Printf("%s", err)
		}

		err = front.ServeTemplate(writer, "pool", users)

		if err != nil {
			log.Printf("%s", err)
		}
	}

	// Read target handle to operate on
	target, err := front.ReadFormStringWithFailure(
		"handle", true, &request.Form, serveStatus, "Target username required!",
	)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	action, err := front.ReadFormRadio(
		"action", []string{"add", "block"}, &request.Form,
		serveStatus, "Action required!",
	)

	// Load pool from current user
	pool, err := data.LoadPool(account.Handle)

	if err != nil {
		log.Printf("Pool not found for user \"%s\"! Rebuilding...", account.Handle)

		pool, err = data.AddPool(account.Handle)

		if err != nil {
			front.Error501(writer)
			log.Printf("%s", err)
			return
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
		serveStatus("User does not exist!")
		log.Printf("%s", err)
		return
	}

	serveStatus("")
}
