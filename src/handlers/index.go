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
		views := control.MakeViews(nil, active)
		return front.ServeTemplate(out, "index", views)
	}

	// Get the feed view for the current user
	view, err := control.MakeFeedView(active)

	if err != nil {
		log.Printf("%s", err)
	}

	views := control.MakeViews(view, active)
	return front.ServeTemplate(out, "index", views)
}
