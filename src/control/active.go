package control

import (
	"../data"
	"../views"
	"bytes"
	"os/exec"
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
	var hash string
	var link string

	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	output := out.String()
	if err != nil || strings.Contains(output, "fatal") {
		hash = "error"
		link = "#"
	} else {
		hash = strings.TrimSpace(output)
		link = "https://github.com/acgaudette/social-humans/commits/" + hash
	}
	return views.Base{
		Commit: hash,
		Link:   link,
	}
}

// Build a Status view
func MakeStatusView(status string) views.Status {
	return views.Status{
		Status: status,
	}
}
