package data

import (
	"../../smhb"
	"net/http"
)

func AddUser(handle, password, name string) (smhb.User, error) {
	return nil, nil
}

func LoadUser(handle string) (smhb.User, error) {
	return nil, nil
}

func RemoveUser(handle string) error {
	return nil
}

func AddPool(handle string) (smhb.Pool, error) {
	return nil, nil
}

func LoadPool(handle string) (smhb.Pool, error) {
	return nil, nil
}

func AddPost(title, content string, author smhb.User) error {
	return nil
}

func LoadPost(address string) (smhb.Post, error) {
	return nil, nil
}

func RemovePost(address string) error {
	return nil
}

func BuildPostAddress(handle, stamp string) string {
	return ""
}

func GetPostAddresses(author string) ([]string, error) {
	return nil, nil
}

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
