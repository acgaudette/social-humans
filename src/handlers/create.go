package handlers

import (
	"../data"
	"../front"
	"log"
	"net/http"
)

func GetCreate(writer http.ResponseWriter, request *http.Request) {
	err := front.ServeTemplate(writer, "create", &front.StatusMessage{})

	if err != nil {
		log.Printf("%s", err)
	}
}

func Create(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	// Serve back the page with a status message
	serveStatus := func(status string) {
		message := front.StatusMessage{Status: status}
		if err := front.ServeTemplate(writer, "create", &message); err != nil {
			log.Printf("%s", err)
		}
	}

	handle, err := front.ReadFormStringWithFailure(
		"handle", true, &request.Form, serveStatus, "Username required!",
	)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	name, err := front.ReadFormStringWithFailure(
		"name", false, &request.Form, serveStatus, "Name required!",
	)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	password, err := front.ReadFormStringWithFailure(
		"password", false, &request.Form, serveStatus, "Password required!",
	)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	// Check for existing user
	account, err := data.LoadUser(handle)

	// If user exists, fail
	if err == nil {
		serveStatus("Username taken!")
		return
	}

	// Add new user
	account, err = data.AddUser(handle, password, name)

	if err != nil {
		front.Error501(writer)
		log.Printf("%s", err)
		return
	}

	// Create session and redirect back home
	data.AddSession(writer, account)
	http.Redirect(writer, request, "/", http.StatusFound)
}
