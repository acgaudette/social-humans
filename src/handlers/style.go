package handlers

import (
	"../app"
	"../control"
	"net/http"
)

func GetStyle(out http.ResponseWriter, in *http.Request) *app.Error {
	http.ServeFile(out, in, control.ROOT+"/style.css")
	return nil
}
