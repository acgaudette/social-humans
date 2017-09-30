package handlers

import (
	"../front"
	"errors"
	"log"
	"net/http"
	"strings"
)

func serveStatusError(template, status string, writer http.ResponseWriter) {
	message := front.StatusMessage{Status: status}
	err := front.ServeTemplate(writer, template, &message)

	if err != nil {
		front.Error501(writer)
		log.Printf("%s", err)
	}
}

func readFormString(
	template, key, errorStatus string,
	writer http.ResponseWriter, request *http.Request,
) (string, error) {
	result := request.Form.Get(key)

	if result == "" {
		serveStatusError(template, errorStatus, writer)
		return "", errors.New("key not found for string")
	}

	if strings.IndexRune(result, '+') >= 0 {
		serveStatusError(template, "Invalid input", writer)
		return "", errors.New("invalid input")
	}

	return result, nil
}
