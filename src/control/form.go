package control

import (
	"../../smhb"
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

// Helper function used by create
func ReadCreateForm(form *url.Values) (*string, *string, *string, *string) {
	handle, err := SanitizeFormString("handle", form)

	if err != nil {
		status := "Invalid username"
		return nil, nil, nil, &status
	}

	if handle == "" {
		status := "Username required!"
		return nil, nil, nil, &status
	}

	// Check handle character limit
	if utf8.RuneCountInString(handle) > smhb.HANDLE_LIMIT {
		status := fmt.Sprintf(
			"Username must be under %v characters", smhb.HANDLE_LIMIT,
		)

		return nil, nil, nil, &status
	}

	name := form.Get("name")

	if name == "" {
		status := "Name required!"
		return &handle, nil, nil, &status
	}

	password := form.Get("password")

	if password == "" {
		status := "Password required!"
		return &handle, &name, nil, &status
	}

	return &handle, &name, &password, nil
}

// Helper function shared by create_post and edit_post
func ReadPostForm(form *url.Values) (*string, *string, *string) {
	title := form.Get("title")

	if title == "" {
		status := "Title required"
		return nil, nil, &status
	}

	// Check title character limit
	if utf8.RuneCountInString(title) > smhb.TITLE_LIMIT {
		status := fmt.Sprintf(
			"Post title must be under %v characters", smhb.TITLE_LIMIT,
		)

		return nil, nil, &status
	}

	content := form.Get("content")

	if content == "" {
		status := "Post content required!"
		return nil, nil, &status
	}

	// Check content character limit
	if utf8.RuneCountInString(content) > smhb.CONTENT_LIMIT {
		status := fmt.Sprintf(
			"Post content must be under %v characters", smhb.CONTENT_LIMIT,
		)

		return nil, nil, &status
	}

	return &title, &content, nil
}
