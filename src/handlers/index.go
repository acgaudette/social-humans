package handlers

import (
	"../control"
	"../data"
	"../front"
	"log"
	"net/http"
)

func Index(writer http.ResponseWriter, request *http.Request) {
	// Load current user, if available
	account, err := data.GetUserFromSession(request)

	// Redirect to login page if there is no session open
	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusFound)
		log.Printf("%s", err)
		return
	}

	err = front.ServeTemplate(
		writer, "index", control.GetUserView(account, request),
	)

	if err != nil {
		log.Printf("%s", err)
	}
}
