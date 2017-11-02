package control

import (
	"../../smhb"
	"../views"
)

/*
	Model to view functions never return nil so that something is always
	rendered
*/

// Build a User view
func MakeUserView(user smhb.User, isActive bool) views.User {
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

	view := views.User{
		Handle:       handle,
		Name:         name,
		IsActiveUser: isActive,
	}

	return view
}
