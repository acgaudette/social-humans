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

func MakeUserView(
	user *data.User, status string, active *data.User,
) (*front.UserView, *front.StatusView) {
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

	view := &front.UserView{
		Handle:       handle,
		Name:         name,
		IsActiveUser: isActive,
	}

	return view, MakeStatusView(status)
}

// Load the active user and build a UserView
func GetUserAndMakeUserView(
	user *data.User, status string, in *http.Request,
) (*front.UserView, *front.StatusView, *data.User) {
	active, _ := data.GetUserFromSession(in)
	view, statusView := MakeUserView(user, status, active)
	return view, statusView, active
}
