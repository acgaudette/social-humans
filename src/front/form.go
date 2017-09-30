package front

import (
	"errors"
	"net/http"
	"strings"
)

type errorClosure func(string)

func ReadFormString(
	template, key, errorStatus string,
	fail errorClosure, request *http.Request,
) (string, error) {
	result := request.Form.Get(key)

	if result == "" {
		fail(errorStatus)
		return "", errors.New("key not found for string")
	}

	if strings.IndexRune(result, '+') >= 0 {
		fail(errorStatus)
		return "", errors.New("invalid input")
	}

	return result, nil
}
