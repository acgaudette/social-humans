package main

import (
	"errors"
	"log"
	"net/http"
)

type router struct {
	mux *http.ServeMux
}

func newRouter() *router {
	this := &router{}
	this.mux = http.NewServeMux()

	this.mux.HandleFunc("/", index)
	this.mux.HandleFunc("/login", login)
	this.mux.HandleFunc("/logout", logout)
	this.mux.HandleFunc("/pool", managePool)

	return this
}

func error501(writer http.ResponseWriter) {
	http.Error(
		writer,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}

func index(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		if request.URL.Path == "/" {

			data, err := getUserFromSession(request)

			if err != nil {
				log.Printf("%s", err)
				http.Redirect(writer, request, "/login", http.StatusFound)
				break
			}

			if err = serveTemplate(writer, "/index.html", data); err != nil {
				log.Printf("%s", err)
			}
		} else {
			http.NotFound(writer, request)
		}
	}
}

func login(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		err := serveTemplate(writer, "/login.html", &statusMessage{Status: ""})

		if err != nil {
			log.Printf("%s", err)
			error501(writer)
		}

	case "POST":
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
}

func logout(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		http.Redirect(writer, request, "/", http.StatusFound)

	case "POST":
		clearSession(writer)
		http.Redirect(writer, request, "/", http.StatusFound)
	}
}

func managePool(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		_, err := getUserFromSession(request)

		if err != nil {
			log.Printf("%s", err)
			http.Redirect(writer, request, "/login", http.StatusFound)
			return
		}

		err = serveTemplate(writer, "/pool.html", statusMessage{Status: ""})

		if err != nil {
			log.Printf("%s", err)
			error501(writer)
		}

	case "POST":
		account, err := getUserFromSession(request)

		if err != nil {
			log.Printf("%s", err)
			http.Redirect(writer, request, "/login", http.StatusFound)
			return
		}

		request.ParseForm()

		serveError := func(status string) error {
			message := statusMessage{Status: status}
			err := serveTemplate(writer, "/pool.html", &message)

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
			loadPoolAndAdd(account.Handle, target)

		case "block":
			loadPoolAndBlock(account.Handle, target)
		}

		http.Redirect(writer, request, "/", http.StatusFound)
	}
}
