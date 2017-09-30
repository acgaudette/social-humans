package handlers

import (
	"../data"
	"../front"
	"log"
	"net/http"
)

func Index(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/" {
		data, err := data.GetUserFromSession(request)

		if err != nil {
			log.Printf("%s", err)
			http.Redirect(writer, request, "/login", http.StatusFound)
			return
		}

		if err = front.ServeTemplate(writer, "/index.html", data); err != nil {
			log.Printf("%s", err)
		}
	} else {
		http.NotFound(writer, request)
	}
}
