package control

import (
	"fmt"
	"net/url"
	"strings"
)

func SanitizeFormString(key string, form *url.Values) (string, error) {
	result := form.Get(key)

	// Look for our delimiter in the input
	if strings.IndexRune(result, '+') >= 0 {
		return "", fmt.Errorf("invalid form input for key \"%s\"", key)
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

	return "", fmt.Errorf("key \"%s\" not found for radio form")
}