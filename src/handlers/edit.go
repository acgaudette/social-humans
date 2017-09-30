package handlers

import (
	"../control"
	"../data"
	"../front"
	"log"
	"net/http"
)

func GetEdit(writer http.ResponseWriter, request *http.Request) {
	account, err := data.GetUserFromSession(request)

	if err != nil {
		log.Printf("%s", err)
		http.Redirect(writer, request, "/login", http.StatusFound)
		return
	}

	err = front.ServeTemplate(
		writer, "edit", control.GetUserView(account, request),
	)

	if err != nil {
		log.Printf("%s", err)
	}
}

func Edit(writer http.ResponseWriter, request *http.Request) {
	account, err := data.GetUserFromSession(request)

	if err != nil {
		log.Printf("%s", err)
		http.Redirect(writer, request, "/login", http.StatusFound)
		return
	}

	request.ParseForm()

	serveError := func(status string) {
		message := control.GetUserView(account, request)
		message.Status = status

		err := front.ServeTemplate(writer, "edit", &message)

		if err != nil {
			log.Printf("%s", err)
		}
	}

	name, err := front.ReadFormString(
		"edit", "name", "", func(string) {}, request,
	)

	if err != nil {
		log.Printf("%s", err)
	}

	if name != "" {
		if err = account.SetName(name); err != nil {
			front.Error501(writer)
			log.Printf("%s", err)
			return
		}
	}

	password, err := front.ReadFormString(
		"create", "password_0", "", func(string) {}, request,
	)

	if err != nil {
		log.Printf("%s", err)
	}

	confirm, err := front.ReadFormString(
		"create", "password_1", "", func(string) {}, request,
	)

	if err != nil {
		log.Printf("%s", err)
	}

	if password == confirm && password != "" {
		if err = account.SetPassword(password); err != nil {
			front.Error501(writer)
			log.Printf("%s", err)
			return
		}
	} else if password != confirm {
		serveError("Passwords don't match!")
		return
	}

	serveError("Information updated")
}
