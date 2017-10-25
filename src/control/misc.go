package control

import (
	"../../smhb"
	"../views"
)

/*
	Model to view functions never return nil so that something is always
	rendered
*/

// Build an Active view from a user
func MakeActiveView(active smhb.User) views.Active {
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
