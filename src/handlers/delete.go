package handlers

import (
	"../data"
	"log"
	"net/http"
)

func GetDelete(writer http.ResponseWriter, request *http.Request) {
	Me(writer, request)
}

func Delete(writer http.ResponseWriter, request *http.Request) {
	// Load current user, if available
	account, err := data.GetUserFromSession(request)

	// Redirect to login page if there is no session open
	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusFound)
		log.Printf("%s", err)
		return
	}

	// Remove user and clear session
	data.RemoveUser(account.Handle)
	Logout(writer, request)
}
