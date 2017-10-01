package handlers

import (
	"../control"
	"../data"
	"../front"
	"log"
	"net/http"
	"strings"
)

func GetUser(writer http.ResponseWriter, request *http.Request) {
	// Extract the handle and attempt to load user
	handle := request.URL.Path[strings.LastIndex(request.URL.Path, "/")+1:]
	account, err := data.LoadUser(handle)

	// User does not exist
	if err != nil {
		http.NotFound(writer, request)
		log.Printf("%s", err)
		return
	}

	err = front.ServeTemplate(
		writer, "user", control.GetUserView(account, request),
	)

	if err != nil {
		log.Printf("%s", err)
	}
}
