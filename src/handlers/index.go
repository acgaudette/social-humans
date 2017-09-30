package handlers

import (
	"../data"
	"../control"
	"../front"
	"log"
	"net/http"
)

func Index(writer http.ResponseWriter, request *http.Request) {
	account, err := data.GetUserFromSession(request)

	if err != nil {
		log.Printf("%s", err)
		http.Redirect(writer, request, "/login", http.StatusFound)
		return
	}

	err = front.ServeTemplate(
		writer, "/index.html", control.GetUserView(account),
	)

	if err != nil {
		log.Printf("%s", err)
	}
}
