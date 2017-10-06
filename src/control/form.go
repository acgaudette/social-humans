package control

import (
	"../app"
	"../data"
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

// Read string from form and check for invalid input
func SanitizeFormString(key string, form *url.Values) (string, error) {
	result := form.Get(key)

	// Look for our delimiter in the input
	if strings.IndexRune(result, '+') >= 0 {
		return "", fmt.Errorf("invalid form input for key \"%s\"", key)
	}

	return result, nil
}

// Read selection from radio
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

// Helper function shared by create_post and edit_post
func ReadPostForm(
	serve func(string) *app.Error, form *url.Values,
) (string, string, *app.Error) {
	title := form.Get("title")

	if title == "" {
		return "", "", serve("Title required!")
	}

	// Check title character limit
	if utf8.RuneCountInString(title) > data.TITLE_LIMIT {
		return "", "", serve(
			fmt.Sprintf(
				"Post title must be under %v characters", data.TITLE_LIMIT,
			),
		)
	}

	content := form.Get("content")

	if content == "" {
		return "", "", serve("Post content required!")
	}

	// Check content character limit
	if utf8.RuneCountInString(content) > data.CONTENT_LIMIT {
		return "", "", serve(
			fmt.Sprintf(
				"Post content must be under %v characters", data.CONTENT_LIMIT,
			),
		)
	}

	return title, content, nil
}
