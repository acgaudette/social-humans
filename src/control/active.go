package control

import (
	"../data"
	"../views"
)

/*
	Model to view functions never return nil, so that something is always
	rendered
*/

// Build a views map from a generic view, a status view, and the active user
func MakeContainer(
	view interface{}, status *views.Status, active *data.User,
) *views.Container {
	views := make(views.Container)

	// Content (main) view
	if view != nil {
		views["content"] = view
	}

	// Active user account view
	if active != nil {
		views["active"] = MakeActiveView(active.Handle)
	}

	// Status message (info) view
	if status != nil {
		views["status"] = status
	}

	return &views
}

// Build an Active view
func MakeActiveView(handle string) *views.Active {
	return &views.Active{
		Handle: handle,
	}
}

// Build a Status view
func MakeStatusView(status string) *views.Status {
	return &views.Status{
		Status: status,
	}
}
