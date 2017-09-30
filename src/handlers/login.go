package handlers

import (
	"../data"
	"../front"
	"log"
	"net/http"
)

func GetLogin(writer http.ResponseWriter, request *http.Request) {
	err := front.ServeTemplate(
		writer, "login", &front.StatusMessage{},
	)

	if err != nil {
		log.Printf("%s", err)
		front.Error501(writer)
	}
}

func Login(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	serveError := func(status string) {
		message := front.StatusMessage{Status: status}
		err := front.ServeTemplate(writer, "login", &message)

		if err != nil {
			front.Error501(writer)
			log.Printf("%s", err)
		}
	}

	handle, err := front.ReadFormString(
		"login", "handle", "Username required!",
		serveError, request,
	)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	password, err := front.ReadFormString(
		"login", "password", "Password required!",
		serveError, request,
	)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	account, err := data.LoadUser(handle)

	if err != nil {
		serveError("User does not exist!")
		return
	} else if err = account.Validate(password); err != nil {
		log.Printf("%s", err)

		message := front.StatusMessage{Status: "Invalid password"}
		err := front.ServeTemplate(writer, "login", &message)

		if err != nil {
			log.Printf("%s", err)
			front.Error501(writer)
		}

		return
	}

	data.AddSession(writer, account)
	http.Redirect(writer, request, "/", http.StatusFound)
}
