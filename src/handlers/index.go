package handlers

import (
	"../../smhb"
	"../app"
	"../control"
	"../data"
	"log"
	"net/http"
)

func Index(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, err := data.GetUserFromSession(in)

	if err != nil {
		// Connection error
		if _, ok := err.(smhb.NotFoundError); !ok {
			return app.ServerError(err)
		}

		// Serve blank index if there is no session open
		container := control.MakeContainer()
		return app.ServeTemplate(out, "index", container)
	}

	// Initialize an empty status view
	status := control.MakeStatusView("")

	// Get the feed view for the current user
	view, err := control.MakeFeedView(active.Handle())

	if err != nil {
		// Update status message with regards to the error
		switch err.(type) {
		case *control.AccessError:
			status.Status = "Error: access failure"
			log.Printf("%s", err)

		case *control.EmptyFeedError:
			status.Status = "Nothing to see here..."

		default:
			log.Printf("%s", err)
		}
	}

	// Build views
	container := control.MakeContainer()
	container.SetActive(control.MakeActiveView(active))
	container.SetContent(view)
	container.SetStatus(status)

	// Serve
	return app.ServeTemplate(out, "index", container)
}
