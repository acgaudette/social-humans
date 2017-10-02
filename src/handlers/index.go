package handlers

import (
	"../app"
	"../control"
	"../front"
	"net/http"
)

func Index(out http.ResponseWriter, in *http.Request) *app.Error {
	views := control.GetUserAndMakeViews(control.MakeFeedView(), in)
	return front.ServeTemplate(out, "index", views)
}
