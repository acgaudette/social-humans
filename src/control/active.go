package control

import (
	"../data"
	"../views"
	"io/ioutil"
	"log"
	"strings"
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
	// Initialize empty view
	view := views.Base{
		Link: "#",
	}

	// Read HEAD
	buffer, err := ioutil.ReadFile(GIT_DIR + "/HEAD")

	if err != nil {
		log.Printf("%s", err)
		// Always display something to the frontend
		return view
	}

	// Convert to string and trim newline
	path := string(buffer[:])
	path = path[:len(path)-1]

	// Get path from HEAD ref
	ref := strings.Split(path, " ")[1]

	// Read commit hash
	buffer, err = ioutil.ReadFile(GIT_DIR + "/" + ref)

	if err != nil {
		log.Printf("%s", err)
		// Always display something to the frontend
		return view
	}

	// Convert to string
	hash := string(buffer[:])

	// Make short hash
	hash = hash[:7]

	// Build view
	view.Commit = hash
	view.Link = "https://github.com/acgaudette/social-humans/commits/" + hash
	return view
}

// Build a Status view
func MakeStatusView(status string) views.Status {
	return views.Status{
		Status: status,
	}
}
