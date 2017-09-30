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
	handle := request.URL.Path[strings.LastIndex(request.URL.Path, "/")+1:]
	account, err := data.LoadUser(handle)

	if err != nil {
		log.Printf("%s", err)
		http.NotFound(writer, request)
		return
	}

	err = front.ServeTemplate(
		writer, "user", control.GetUserView(account, request),
	)

	if err != nil {
		log.Printf("%s", err)
	}
}
