package data

import (
	"../../smhb"
	"net/http"
)

func AddSession(out http.ResponseWriter, account smhb.User) error {
	return nil
}

func JoinSession(out http.ResponseWriter, account smhb.User) error {
	return nil
}

func ClearSession(out http.ResponseWriter) {}

func GetUserFromSession(in *http.Request) (smhb.User, error) {
	return nil, nil
}
