package front

import (
	"net/http"
)

func Error403(writer http.ResponseWriter) {
	http.Error(
		writer,
		http.StatusText(http.StatusForbidden),
		http.StatusForbidden,
	)
}

func Error501(writer http.ResponseWriter) {
	http.Error(
		writer,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}
