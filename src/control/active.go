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
	container := views.NewContainer()
	container.SetBase(MakeBaseView())
	return container
}

// Build an Active view from a user
func MakeActiveView(active data.User) views.Active {
	return views.Active{
		Handle: active.Handle(),
	}
}

// Build a Base view
func MakeBaseView() views.Base {
	// Check if commit hash has been loaded
	if data.CommitHash != nil {
		prefix := "https://github.com/acgaudette/social-humans/commits/"

		return views.Base{
			Commit: *data.CommitHash,
			Link:   prefix + *data.CommitHash,
		}
	}

	// Otherwise, return empty view
	return views.Base{
		Link: "#",
	}
}

// Build a Status view
func MakeStatusView(status string) views.Status {
	return views.Status{
		Status: status,
	}
}
