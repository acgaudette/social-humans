package control

import (
	"../data"
	"../front"
)

/*
	Model to view functions never return nil, so that something is always
	rendered
*/

// Build a views map from a generic view, a status view, and the active user
func MakeViews(
	view interface{}, status *front.StatusView, active *data.User,
) *front.Views {
	views := make(front.Views)

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
