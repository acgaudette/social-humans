package front

import (
	"../app"
	"net/http"
)

// Redirect client and return error for router to handle
func Redirect(
	route string, err error, out http.ResponseWriter, in *http.Request,
) *app.Error {
	http.Redirect(out, in, route, http.StatusFound)

	return &app.Error{
		Native: err,
		Code:   app.REDIRECT,
	}
}

// Return server error (501) for router to handle
func ServerError(err error) *app.Error {
	return &app.Error{
		Native: err,
		Code:   app.SERVER,
	}
}
