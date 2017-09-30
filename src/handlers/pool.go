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
			log.Printf("%s", err)
		}
	}

	target, err := front.ReadFormString(
		"pool", "handle", "Target username required!",
		serveError, request,
	)

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

	pool, err := data.LoadPool(account.Handle)

	if err != nil {
		log.Printf("%s", err)
		log.Printf("Pool not found for user \"%s\"! Rebuilding...", account.Handle)

		pool, err = data.AddPool(account.Handle)

		if err != nil {
			log.Printf("%s", err)
			return
		}
	}

	switch action {
	case "add":
		err = pool.Add(target)

	case "block":
		err = pool.Block(target)
	}

	if err != nil {
		serveError("User does not exist!")
		log.Printf("%s", err)
		return
	}

	serveError("")
}
