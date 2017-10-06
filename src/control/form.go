package control

import (
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
func ReadPostForm(form *url.Values) (*string, *string, *string) {
	title := form.Get("title")

	if title == "" {
		status := "Title required"
		return nil, nil, &status
	}

	// Check title character limit
	if utf8.RuneCountInString(title) > data.TITLE_LIMIT {
		status := fmt.Sprintf(
			"Post title must be under %v characters", data.TITLE_LIMIT,
		)

		return nil, nil, &status
	}

	content := form.Get("content")

	if content == "" {
		status := "Post content required!"
		return nil, nil, &status
	}

	// Check content character limit
	if utf8.RuneCountInString(content) > data.CONTENT_LIMIT {
		status := fmt.Sprintf(
			"Post content must be under %v characters", data.CONTENT_LIMIT,
		)

		return nil, nil, &status
	}

	return &title, &content, nil
}
