package front

import (
	"../app"
	"net/http"
)

func Error403(out http.ResponseWriter) {
	http.Error(
		out,
		http.StatusText(http.StatusForbidden),
		http.StatusForbidden,
	)
}

func Error501(out http.ResponseWriter) {
	http.Error(
		out,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}

func Redirect(
	route string, err error, out http.ResponseWriter, in *http.Request,
) *app.Error {
	http.Redirect(out, in, route, http.StatusFound)

	return &app.Error{
		Native: err,
		Code:   app.REDIRECT,
	}
}
