package handlers

import (
	"../data"
	"../front"
	"errors"
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

func GetLogin(writer http.ResponseWriter, request *http.Request) {
	err := front.ServeTemplate(
		writer, "/login.html", &front.StatusMessage{Status: ""},
	)

	if err != nil {
		log.Printf("%s", err)
		front.Error501(writer)
	}
}

func Login(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()

	readString := func(key string, errorStatus string) (string, error) {
		result := request.Form.Get(key)

		if result == "" {
			message := front.StatusMessage{Status: errorStatus}
			err := front.ServeTemplate(writer, "/login.html", &message)

			if err != nil {
				front.Error501(writer)
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

	account, err := data.LoadUser(handle)

	if err != nil {
		log.Printf("%s", err)

		account, err = data.AddUser(handle, password)

		if err != nil {
			log.Printf("%s", err)
			front.Error501(writer)
			return
		}
	} else if err = account.Validate(password); err != nil {
		log.Printf("%s", err)

		message := front.StatusMessage{Status: "Invalid password"}
		err := front.ServeTemplate(writer, "/login.html", &message)

		if err != nil {
			log.Printf("%s", err)
			front.Error501(writer)
		}

		return
	}

	data.AddSession(writer, account)
	http.Redirect(writer, request, "/", http.StatusFound)
}

func GetLogout(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/", http.StatusFound)
}

func Logout(writer http.ResponseWriter, request *http.Request) {
	data.ClearSession(writer)
	http.Redirect(writer, request, "/", http.StatusFound)
}

func GetPool(writer http.ResponseWriter, request *http.Request) {
	account, err := data.GetUserFromSession(request)

	if err != nil {
		log.Printf("%s", err)
		http.Redirect(writer, request, "/login", http.StatusFound)
		return
	}

	users, err := front.GetPoolUsers(account.Handle, "")

	if err != nil {
		log.Printf("%s", err)
	}

	err = front.ServeTemplate(writer, "/pool.html", users)

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

	serveError := func(status string) error {
		users, err := front.GetPoolUsers(account.Handle, status)

		if err != nil {
			log.Printf("%s", err)
		}

		err = front.ServeTemplate(writer, "/pool.html", users)

		if err != nil {
			front.Error501(writer)
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
		if err = data.LoadPoolAndAdd(account.Handle, target); err != nil {
			log.Printf("%s", err)
		}

	case "block":
		if err = data.LoadPoolAndBlock(account.Handle, target); err != nil {
			log.Printf("%s", err)
		}
	}

	http.Redirect(writer, request, "/pool", http.StatusFound)
}
