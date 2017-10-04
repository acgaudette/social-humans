package handlers

import (
	"../app"
	"net/http"
)

func GetStyle(out http.ResponseWriter, in *http.Request) *app.Error {
	http.ServeFile(out, in, app.ROOT+"/style.css")
	return nil
}
