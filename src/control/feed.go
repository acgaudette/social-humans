package control

import (
	"../data"
	"../views"
	"log"
	"strconv"
	"strings"
)

/*
	Model to view functions never return nil, so that something is always
	rendered
*/

// Build Feed view from a user model
func MakeFeedView(account *data.User) (*views.Feed, error) {
	// Create empty feed
	feed := &views.Feed{
		Posts: []*views.Post{},
	}

	if account == nil {
		// Return empty feed view if user is not found
		return feed, &UserNotFoundError{account.Handle}
	}

	pool, err := data.LoadPool(account.Handle)

	if err != nil {
		log.Printf("error while accessing feed: %s", err)

		// Return empty feed view if pool is not found
		return feed, &AccessError{account.Handle}
	}

	q := FeedQueue{}

	// Iterate through pool and get the user posts
	for _, handle := range pool.Users {
		addresses, err := data.GetPostAddresses(handle)

		if err != nil {
			log.Printf("Error getting posts from \"%s\": %s", handle)
			continue
		}

		// Iterate through posts and push to the priority queue
		for _, post := range addresses {
			score, err := ScorePost(post)

			if err != nil {
				log.Printf("Error parsing address \"%s\": %s", post)
				continue
			}

			q.Add(post, score)
		}
	}

	// Convert queue into feed view
	for q.Len() > 0 {
		address := q.Remove()
		post, err := data.LoadPost(address)

		if err != nil {
			// Always display something to the frontend
			post = &data.Post{
				Title:   "Title Invalid",
				Content: "Content Invalid",
				Author:  "Author Invalid",
			}

			log.Printf("error while updating feed: %s", err)
		}

		// Assumes the account passed in is the active user
		view := MakePostView(post, account)

		feed.Posts = append(feed.Posts, view)
	}

	if len(feed.Posts) == 0 {
		return feed, &EmptyFeedError{account.Handle}
	}

	return feed, nil
}

// Assign a priority to a post
func ScorePost(address string) (int, error) {
	stamp := strings.Split(address, "/")[1]
	result, err := strconv.Atoi(stamp)

	if err != nil {
		return -1, err
	}

	return result, err
}
