package front

import (
	"errors"
	"net/url"
	"strings"
)

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

func ReadFormRadio(
	key string, options []string, form *url.Values,
) (string, error) {
	for _, value := range options {
		if value == form.Get(key) {
			return value, nil
		}
	}

	return "", errors.New("key not found for radio")
}
