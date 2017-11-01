package control

import (
	"../../smhb"
	"../data"
	"../views"
	"log"
)

/*
	Model to view functions never return nil so that something is always
	rendered
*/

// Build Feed view from a user handle
func MakeFeedView(handle string, token smhb.Token) (*views.Feed, error) {
	// Create empty feed
	feed := &views.Feed{
		Posts: []views.Post{},
	}

	loaded, err := data.Backend.GetFeed(handle, token)

	if err != nil {
		log.Printf("error while accessing feed: %s", err)

		// Return empty feed view if pool is not found
		return feed, &AccessError{handle}
	}

	// Iterate through feed addresses, load the associated posts,
	// build the feed view
	for _, address := range loaded.Addresses() {
		post, err := data.Backend.GetPost(handle, address, token)

		if err != nil {
			// Always display something to the frontend
			feed.Posts = append(feed.Posts, emptyPostView())

			log.Printf("error while updating feed: %s", err)
			continue
		}

		// Assumes the account passed in is the active user
		view := MakePostView(post, post.WasAuthoredBy(handle))
		feed.Posts = append(feed.Posts, view)
	}

	if len(feed.Posts) == 0 {
		return feed, &EmptyFeedError{handle}
	}

	return feed, nil
}
