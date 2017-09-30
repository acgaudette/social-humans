package handlers

import (
	"../data"
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

func GetCreate(writer http.ResponseWriter, request *http.Request) {
	err := front.ServeTemplate(
		writer, "create", &front.StatusMessage{},
	)

	if err != nil {
		log.Printf("%s", err)
		front.Error501(writer)
	}
}

func Create(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	handle, err := readFormString(
		"create", "handle", "Username required!", writer, request,
	)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	name, err := readFormString(
		"create", "name", "Name required!", writer, request,
	)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	password, err := readFormString(
		"create", "password", "Password required!", writer, request,
	)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	// Check for existing user
	account, err := data.LoadUser(handle)

	// If user exists, fail
	if err == nil {
		serveStatusError("create", "Username taken!", writer)
		return
	}

	account, err = data.AddUser(handle, password, name)

	if err != nil {
		log.Printf("%s", err)
		front.Error501(writer)
		return
	}

	data.AddSession(writer, account)
	http.Redirect(writer, request, "/", http.StatusFound)
}
