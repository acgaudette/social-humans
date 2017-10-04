package handlers

import (
	"../app"
	"../control"
	"../data"
	"../front"
	"log"
	"net/http"
)

func Index(out http.ResponseWriter, in *http.Request) *app.Error {
	// Load current user, if available
	active, err := data.GetUserFromSession(in)

	// Serve blank index if there is no session open
	if err != nil {
		views := control.MakeViews(nil, nil, active)
		return front.ServeTemplate(out, "index", views)
	}

	// Initialize an empty status view
	status := control.MakeStatusView("")

	// Get the feed view for the current user
	view, err := control.MakeFeedView(active)

	if err != nil {
		// Update status message with regards to the error
		switch err.(type) {
		case *control.UserNotFoundError:
			status.Status = "Error: user not found"
			log.Printf("%s", err)

		case *control.AccessError:
			status.Status = "Error: access failure"
			log.Printf("%s", err)

		case *control.EmptyFeedError:
			status.Status = "Nothing to see here..."

		default:
			log.Printf("%s", err)
		}
	}

	// Build views and serve
	views := control.MakeViews(view, status, active)
	return front.ServeTemplate(out, "index", views)
}
