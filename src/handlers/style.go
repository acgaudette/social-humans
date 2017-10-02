package handlers

import (
	"../app"
	"../front"
	"net/http"
)

func GetStyle(out http.ResponseWriter, in *http.Request) *app.Error {
	http.ServeFile(out, in, front.ROOT + "/style.css")
	return nil
}
