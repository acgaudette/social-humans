package control

import (
	"../data"
	"../views"
	"log"
	"time"
)

/*
	Model to view functions never return nil, so that something is always
	rendered
*/

// Build a Post view from a post model
func MakePostView(post *data.Post, active *data.User) *views.Post {
	// Build timestamp
	timestamp := "unknown time"
	location, err := time.LoadLocation("Local")

	if err != nil {
		log.Printf("error while rendering post: %s", err)
	} else {
		time, err := time.Parse(data.TIMESTAMP_LAYOUT, post.Timestamp)

		if err != nil {
			log.Printf("error while rendering post: %s", err)
		} else {
			timestamp = time.In(location).Format(data.HUMAN_TIME_LAYOUT)
		}
	}

	isActive := false

	// Compare the active user to the post author
	if active != nil && active.Handle == post.Author {
		isActive = true
	}

	return &views.Post{
		Title:        post.Title,
		Content:      post.Content,
		Author:       post.Author,
		ID:           post.Timestamp,
		WasEdited:    post.WasEdited,
		Timestamp:    timestamp,
		IsActiveUser: isActive,
	}
}
