package handlers

import (
	"../control"
	"../data"
	"../front"
	"log"
	"net/http"
)

func GetEdit(writer http.ResponseWriter, request *http.Request) {
	// Load current user, if available
	account, err := data.GetUserFromSession(request)

	// Redirect to login page if there is no session open
	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusFound)
		log.Printf("%s", err)
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

	// Serve back the page with a status message
	serveStatus := func(status string) {
		message := control.GetUserView(account, request)
		message.Status = status

		if err := front.ServeTemplate(writer, "edit", &message); err != nil {
			log.Printf("%s", err)
		}
	}

	request.ParseForm()

	name, err := front.ReadFormString("name", false, &request.Form)

	if err != nil {
		log.Printf("%s", err)
	}

	// Set new user full name
	if name != "" {
		if err = account.SetName(name); err != nil {
			front.Error501(writer)
			log.Printf("%s", err)
			return
		}
	}

	password, err := front.ReadFormString("newPassword", false, &request.Form)

	if err != nil {
		log.Printf("%s", err)
	}

	confirm, err := front.ReadFormString("confirmPassword", false, &request.Form)

	if err != nil {
		log.Printf("%s", err)
	}

	old, err := front.ReadFormString("oldPassword", false, &request.Form)

	if err != nil {
		log.Printf("%s", err)
	}

	// Check if new passwords match and are valid
	if password == confirm && password != "" {
		// Validate password
		if err = account.Validate(old); err != nil {
			serveStatus("Invalid password")
			log.Printf("%s", err)
			return
		}

		// Set new user password
		if err = account.UpdatePassword(password); err != nil {
			front.Error501(writer)
			log.Printf("%s", err)
			return
		}

	} else if password != confirm {
		serveStatus("Passwords don't match!")
		return
	}

	serveStatus("Information updated")
}
