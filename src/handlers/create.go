package handlers

import (
	"../data"
	"../front"
	"log"
	"net/http"
)

func GetCreate(writer http.ResponseWriter, request *http.Request) {
	err := front.ServeTemplate(
		writer, "create", &front.StatusMessage{},
	)

	if err != nil {
		log.Printf("%s", err)
		front.Error501(writer)
	}
}

func Create(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	serveError := func(status string) {
		message := front.StatusMessage{Status: status}
		err := front.ServeTemplate(writer, "create", &message)

		if err != nil {
			front.Error501(writer)
			log.Printf("%s", err)
		}
	}

	handle, err := readFormString(
		"create", "handle", "Username required!",
		serveError, request,
	)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	name, err := readFormString(
		"create", "name", "Name required!",
		serveError, request,
	)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	password, err := readFormString(
		"create", "password", "Password required!",
		serveError, request,
	)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	// Check for existing user
	account, err := data.LoadUser(handle)

	// If user exists, fail
	if err == nil {
		serveError("Username taken!")
		return
	}

	account, err = data.AddUser(handle, password, name)

	if err != nil {
		log.Printf("%s", err)
		front.Error501(writer)
		return
	}

	data.AddSession(writer, account)
	http.Redirect(writer, request, "/", http.StatusFound)
}
