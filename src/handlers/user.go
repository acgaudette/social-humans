package handlers

import (
	"../data"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func GetUser(writer http.ResponseWriter, request *http.Request) {
	handle := request.URL.Path[strings.LastIndex(request.URL.Path, "/")+1:]
	account, err := data.LoadUser(handle)

	if err != nil {
		fmt.Fprintf(writer, "User does not exist!")

		log.Printf("%s", err)
		return
	}

	fmt.Fprintf(writer, "User: "+account.Handle)
}
