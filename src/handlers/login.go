package handlers

import (
	"../data"
	"../front"
	"errors"
	"log"
	"net/http"
)

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
