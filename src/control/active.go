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

// Build a views map from the active user and a generic view
func MakeViews(view interface{}, activeUser *data.User) *front.Views {
	views := make(front.Views)

	// Content (main) view
	if view != nil {
		views["content"] = view
	}

	// Active user account view
	if activeUser != nil {
		views["active"] = MakeActiveView(activeUser.Handle)
	}

	return &views
}

// Build a views map from the active user and a generic view with a status
func MakeViewsWithStatus(
	view interface{}, activeUser *data.User, status *front.StatusView,
) *front.Views {
	views := make(front.Views)

	// Content (main) view
	if view != nil {
		views["content"] = view
	}

	// Active user account view
	if activeUser != nil {
		views["active"] = MakeActiveView(activeUser.Handle)
	}

	views["status"] = status

	return &views
}

// MakeViews, but automatically load the active user
func GetUserAndMakeViews(view interface{}, in *http.Request) *front.Views {
	account, _ := data.GetUserFromSession(in) // Nil check done in MakeViews
	return MakeViews(view, account)
}

// Build an ActiveView
func MakeActiveView(handle string) *front.ActiveView {
	return &front.ActiveView{
		Handle: handle,
	}
}

// Build a StatusView
func MakeStatusView(status string) *front.StatusView {
	return &front.StatusView{
		Status: status,
	}
}
