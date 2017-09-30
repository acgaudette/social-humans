package handlers

import (
	"../data"
	"net/http"
)

func GetLogout(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/", http.StatusFound)
}

func Logout(writer http.ResponseWriter, request *http.Request) {
	data.ClearSession(writer)
	http.Redirect(writer, request, "/", http.StatusFound)
}
