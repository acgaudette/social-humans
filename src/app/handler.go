package app

import (
	"log"
	"net/http"
)

type Handler func(http.ResponseWriter, *http.Request) *Error

func Handle(handler Handler, out http.ResponseWriter, in *http.Request) {
	// Execute handler
	err := handler(out, in)

	// Handle error responses
	if err != nil {
		switch err.Code {
		case SERVER:
			// Error 501
			http.Error(
				out,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)

		case NOT_FOUND:
			http.NotFound(out, in)

		case REDIRECT:
			// Redirect is taken care of in the handler
		}

		// Log error message
		if err.Native != nil {
			log.Printf("%s", err.Native)
		}
	}

	// Otherwise, response was sent in the handler
}
