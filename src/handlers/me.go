package handlers

import (
	"../data"
	"log"
	"net/http"
)

func Me(writer http.ResponseWriter, request *http.Request) {
	account, err := data.GetUserFromSession(request)

	if err != nil {
		log.Printf("%s", err)
		http.Redirect(writer, request, "/login", http.StatusFound)
		return
	}

	http.Redirect(writer, request, "/user/"+account.Handle, http.StatusFound)
}
