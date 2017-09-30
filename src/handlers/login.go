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

	handle, err := readFormString(
		"login", "handle", "Username required!", writer, request,
	)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	password, err := readFormString(
		"login", "password", "Password required!", writer, request,
	)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	account, err := data.LoadUser(handle)

	if err != nil {
		serveStatusError("login", "User does not exist!", writer)
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
