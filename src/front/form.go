package front

import (
	"errors"
	"net/url"
	"strings"
)

type errorClosure func(string)

func ReadFormString(
	key string, sanitize bool, form *url.Values,
) (string, error) {
	result := form.Get(key)

	// Look for our delimiter in the input
	if sanitize && strings.IndexRune(result, '+') >= 0 {
		return "", errors.New("invalid input")
	}

	return result, nil
}

func ReadFormStringWithFailure(
	key string, sanitize bool, form *url.Values,
	fail errorClosure, notFoundMessage string,
) (string, error) {
	result, err := ReadFormString(key, sanitize, form)

	if err != nil {
		fail("Invalid input")
		return "", err
	}

	if result == "" {
		fail(notFoundMessage)
		return "", errors.New("key not found for string")
	}

	return result, nil
}

func ReadFormRadio(
	key string, options []string, form *url.Values,
	fail errorClosure, notFoundMessage string,
) (string, error) {
	for _, value := range options {
		if value == form.Get(key) {
			return value, nil
		}
	}

	fail(notFoundMessage)
	return "", errors.New("key not found for radio")
}
