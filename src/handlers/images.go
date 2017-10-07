package handlers

import (
	"../app"
	"net/http"
	"strings"
)

func GetImage(out http.ResponseWriter, in *http.Request) *app.Error {
	tokens := strings.Split(in.URL.Path, "/")
	filename := tokens[len(tokens)-1]
	http.ServeFile(out, in, app.ROOT+"/images/"+filename)
	return nil
}
