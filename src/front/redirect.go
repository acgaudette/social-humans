package front

import (
	"../app"
	"net/http"
)

func Redirect(
	route string, err error, out http.ResponseWriter, in *http.Request,
) *app.Error {
	http.Redirect(out, in, route, http.StatusFound)

	return &app.Error{
		Native: err,
		Code:   app.REDIRECT,
	}
}
