package handlers

import (
	"../control"
	"../data"
	"../front"
	"errors"
	"log"
	"net/http"
)

func GetPool(writer http.ResponseWriter, request *http.Request) {
	account, err := data.GetUserFromSession(request)

	if err != nil {
		log.Printf("%s", err)
		http.Redirect(writer, request, "/login", http.StatusFound)
		return
	}

	users, err := control.GetPoolView(account.Handle, "")

	if err != nil {
		log.Printf("%s", err)
	}

	err = front.ServeTemplate(writer, "pool", users)

	if err != nil {
		log.Printf("%s", err)
		front.Error501(writer)
	}
}

func ManagePool(writer http.ResponseWriter, request *http.Request) {
	account, err := data.GetUserFromSession(request)

	if err != nil {
		log.Printf("%s", err)
		http.Redirect(writer, request, "/login", http.StatusFound)
		return
	}

	request.ParseForm()

	serveError := func(status string) {
		// 'users' will never be nil
		users, err := control.GetPoolView(account.Handle, status)

		if err != nil {
			log.Printf("%s", err)
		}

		err = front.ServeTemplate(writer, "pool", users)

		if err != nil {
			front.Error501(writer)
			log.Printf("%s", err)
		}
	}

	readString := func(key string, errorStatus string) (string, error) {
		result := request.Form.Get(key)

		if result == "" {
			serveError(errorStatus)
			return "", errors.New("key not found for string")
		}

		return result, nil
	}

	target, err := readString("handle", "Target username required!")

	if err != nil {
		log.Printf("%s", err)
		return
	}

	readRadio := func(
		key string, options []string, errorStatus string,
	) (string, error) {
		for _, value := range options {
			if value == request.Form.Get(key) {
				return value, nil
			}
		}

		serveError(errorStatus)
		return "", errors.New("key not found for radio")
	}

	action, err := readRadio(
		"action", []string{"add", "block"}, "Action required!",
	)

	switch action {
	case "add":
		if err = data.LoadPoolAndAdd(account.Handle, target); err != nil {
			log.Printf("%s", err)
		}

	case "block":
		if err = data.LoadPoolAndBlock(account.Handle, target); err != nil {
			log.Printf("%s", err)
		}
	}

	if err != nil {
		serveError("User does not exist")
	} else {
		serveError("")
	}
}
