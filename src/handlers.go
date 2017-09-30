package main

import (
	"errors"
	"log"
	"net/http"
)

func index(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path == "/" {
		data, err := getUserFromSession(request)

		if err != nil {
			log.Printf("%s", err)
			http.Redirect(writer, request, "/login", http.StatusFound)
			return
		}

		if err = serveTemplate(writer, "/index.html", data); err != nil {
			log.Printf("%s", err)
		}
	} else {
		http.NotFound(writer, request)
	}
}

func getLogin(writer http.ResponseWriter, request *http.Request) {
	err := serveTemplate(writer, "/login.html", &statusMessage{Status: ""})

	if err != nil {
		log.Printf("%s", err)
		error501(writer)
	}
}

func login(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	readString := func(key string, errorStatus string) (string, error) {
		result := request.Form.Get(key)

		if result == "" {
			message := statusMessage{Status: errorStatus}
			err := serveTemplate(writer, "/login.html", &message)

			if err != nil {
				error501(writer)
				return "", err
			}

			return "", errors.New("key not found for string")
		}

		return result, nil
	}

	handle, err := readString("handle", "Username required!")

	if err != nil {
		log.Printf("%s", err)
		return
	}

	password, err := readString("password", "Password required!")

	if err != nil {
		log.Printf("%s", err)
		return
	}

	account, err := loadUser(handle)

	if err != nil {
		log.Printf("%s", err)

		account, err = addUser(handle, password)

		if err != nil {
			log.Printf("%s", err)
			error501(writer)
			return
		}
	} else if err = account.validate(password); err != nil {
		log.Printf("%s", err)

		message := statusMessage{Status: "Invalid password"}
		err := serveTemplate(writer, "/login.html", &message)

		if err != nil {
			log.Printf("%s", err)
			error501(writer)
		}

		return
	}

	addSession(writer, account)
	http.Redirect(writer, request, "/", http.StatusFound)
}

func getLogout(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/", http.StatusFound)
}

func logout(writer http.ResponseWriter, request *http.Request) {
	clearSession(writer)
	http.Redirect(writer, request, "/", http.StatusFound)
}

func getPool(writer http.ResponseWriter, request *http.Request) {
	account, err := getUserFromSession(request)

	if err != nil {
		log.Printf("%s", err)
		http.Redirect(writer, request, "/login", http.StatusFound)
		return
	}

	users, err := getPoolUsers(account.Handle, "")

	if err != nil {
		log.Printf("%s", err)
	}

	err = serveTemplate(writer, "/pool.html", users)

	if err != nil {
		log.Printf("%s", err)
		error501(writer)
	}
}

func managePool(writer http.ResponseWriter, request *http.Request) {
	account, err := getUserFromSession(request)

	if err != nil {
		log.Printf("%s", err)
		http.Redirect(writer, request, "/login", http.StatusFound)
		return
	}

	request.ParseForm()

	serveError := func(status string) error {
		users, err := getPoolUsers(account.Handle, status)

		if err != nil {
			log.Printf("%s", err)
		}

		err = serveTemplate(writer, "/pool.html", users)

		if err != nil {
			error501(writer)
		}

		return err
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
		if err = loadPoolAndAdd(account.Handle, target); err != nil {
			log.Printf("%s", err)
		}

	case "block":
		if err = loadPoolAndBlock(account.Handle, target); err != nil {
			log.Printf("%s", err)
		}
	}

	http.Redirect(writer, request, "/pool", http.StatusFound)
}
