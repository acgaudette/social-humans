package control

import (
	"../data"
	"../views"
)

/*
	Model to view functions never return nil, so that something is always
	rendered
*/

// Build a view container from the active user
func MakeContainer(active *data.User) *views.Container {
	container := views.NewContainer()

	// Active user account view
	if active != nil {
		container.SetActive(MakeActiveView(active.Handle))
	}

	return &container
}

// Build an Active view
func MakeActiveView(handle string) views.Active {
	return views.Active{
		Handle: handle,
	}
}

// Build a Status view
func MakeStatusView(status string) views.Status {
	return views.Status{
		Status: status,
	}
}
