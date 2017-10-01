package handlers

import (
	"../data"
	"log"
	"net/http"
)

func Me(writer http.ResponseWriter, request *http.Request) {
	// Load current user, if available
	account, err := data.GetUserFromSession(request)

	// Redirect to login page if there is no session open
	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusFound)
		log.Printf("%s", err)
		return
	}

	// Redirect to user page
	http.Redirect(writer, request, "/user/"+account.Handle, http.StatusFound)
}
