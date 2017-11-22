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

	// Set base view automatically
	container.SetBase(MakeBaseView())

	return container
}

// Build a Base view
func MakeBaseView() views.Base {
	// Check if commit hash has been loaded
	if data.CommitHash != nil {
		prefix := "https://github.com/acgaudette/social-humans/commits/"

		return views.Base{
			Commit: *data.CommitHash, // Safe
			Link:   prefix + *data.CommitHash,
		}
	}

	// Otherwise, return empty view
	return views.Base{
		Link: "#",
	}
}
