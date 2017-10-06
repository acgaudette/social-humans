package app

import (
	"net/http"
)

// Redirect client and return error for router to handle
func Redirect(
	route string, err error, out http.ResponseWriter, in *http.Request,
) *Error {
	http.Redirect(out, in, route, http.StatusFound)

	return &Error{
		Native: err,
		Code:   REDIRECT,
	}
}

// Return server error (501) for router to handle
func ServerError(err error) *Error {
	return &Error{
		Native: err,
		Code:   SERVER,
	}
}

// Return not found error (404) for router to handle
func NotFound(err error) *Error {
	return &Error{
		Native: err,
		Code:   NOT_FOUND,
	}
}

// Return forbidden error (403) for router to handle
func Forbidden(err error) *Error {
	return &Error{
		Native: err,
		Code:   FORBIDDEN,
	}
}
