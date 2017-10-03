package control

import (
	"../data"
	"../front"
	"net/http"
)

/*
	Model to view functions never return nil, so that something is always
	rendered
*/

// Build a UserView from a user model
func MakeUserView(
	user *data.User, status string, active *data.User,
) *front.UserView {
	handle := user.Handle

	// Always display something to the frontend
	if handle == "" {
		handle = "Username Invalid"
	}

	name := user.Name

	// Always display something to the frontend
	if name == "" {
		name = "Name Invalid"
	}

	// Compare the active user to the input user
	isActive := false

	if active != nil && active.Handle == user.Handle {
		isActive = true
	}

	return &front.UserView{
		Handle:       handle,
		Name:         name,
		Status:       status,
		IsActiveUser: isActive,
	}
}

// Load the active user and build a UserView
func GetUserAndMakeUserView(
	user *data.User, status string, in *http.Request,
) (*front.UserView, *data.User) {
	active, _ := data.GetUserFromSession(in)
	return MakeUserView(user, status, active), active
}
