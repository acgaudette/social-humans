package control

import (
	"../data"
	"../views"
)

/*
	Model to view functions never return nil so that something is always
	rendered
*/

// Build a User view
func MakeUserView(user data.User, active data.User) *views.User {
	handle := user.Handle()

	// Always display something to the frontend
	if handle == "" {
		handle = "Username Invalid"
	}

	name := user.Name()

	// Always display something to the frontend
	if name == "" {
		name = "Name Invalid"
	}

	// Compare the active user to the input user
	isActive := false

	if active != nil && active.Handle() == user.Handle() {
		isActive = true
	}

	view := &views.User{
		Handle:       handle,
		Name:         name,
		IsActiveUser: isActive,
	}

	return view
}
