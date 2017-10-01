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
	}
}

func Login(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	serveStatus := func(status string) {
		message := front.StatusMessage{Status: status}
		if err := front.ServeTemplate(writer, "login", &message); err != nil {
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

	password, err := front.ReadFormStringWithFailure(
		"password", false, &request.Form, serveStatus, "Password required!",
	)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	// Load user account
	account, err := data.LoadUser(handle)

	if err != nil {
		serveStatus("User does not exist!")
		return

		// Validate password
	} else if err = account.Validate(password); err != nil {
		serveStatus("Invalid password")
		log.Printf("%s", err)
		return
	}

	// Create user session and redirect back home
	data.AddSession(writer, account)
	http.Redirect(writer, request, "/", http.StatusFound)
}
