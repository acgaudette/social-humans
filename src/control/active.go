package control

import (
	"../data"
	"../views"
)

/*
	Model to view functions never return nil so that something is always
	rendered
*/

// Build a view container
func MakeContainer() views.Container {
	return views.NewContainer()
}

// Build an Active view from a user
func MakeActiveView(active data.User) views.Active {
	return views.Active{
		Handle: active.Handle(),
	}
}

// Build a Status view
func MakeStatusView(status string) views.Status {
	return views.Status{
		Status: status,
	}
}
