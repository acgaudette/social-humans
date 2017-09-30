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
	account, err := data.GetUserFromSession(request)

	if err != nil {
		log.Printf("%s", err)
		http.Redirect(writer, request, "/login", http.StatusFound)
		return
	}

	data.RemoveUser(account.Handle)
	Logout(writer, request)
}
