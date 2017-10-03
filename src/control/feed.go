package control

import (
	"../data"
	"../front"
	"log"
	"net/http"
	"strconv"
	"strings"
)

/*
	Model to view functions never return nil, so that something is always
	rendered
*/

// Build FeedView from a user model
func MakeFeedView(account *data.User) (*front.FeedView, error) {
	// Create empty feed
	feed := &front.FeedView{
		Posts: []*front.PostView{},
	}

	if account == nil {
		// Return empty feed view if user is not found
		feed.Status = "Error: user not found"
		return feed, nil
	}

	pool, err := data.LoadPool(account.Handle)

	if err != nil {
		// Return empty feed view if pool is not found
		feed.Status = "Error: access failure"
		return feed, err
	}

	q := PQueue{}

	// Iterate through pool and push posts to the priority queue
	for _, handle := range pool.Users {
		addresses, err := data.GetPostAddresses(handle)

		if err != nil {
			log.Printf("Error getting posts from \"%s\"", handle)
			continue
		}

		for _, post := range addresses {
			q.Add(post, ScorePost(post))
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

			log.Printf("%s", err)
		}

		// Assumes the account passed in is the active user
		view := MakePostView(post, account)

		feed.Posts = append(feed.Posts, view)
	}

	if len(feed.Posts) == 0 {
		feed.Status = "Nothing to see here..."
	}

	return feed, nil
}

// Build a PostView from a post model
func MakePostView(post *data.Post, active *data.User) *front.PostView {
	isActive := false

	// Compare the active user to the post author
	if active != nil && active.Handle == post.Author {
		isActive = true
	}

	return &front.PostView{
		Title:        post.Title,
		Content:      post.Content,
		Author:       post.Author,
		ID:           post.Timestamp,
		IsActiveUser: isActive,
	}
}

func GetUserAndMakePostView(
	post *data.Post, in *http.Request,
) (*front.PostView, *data.User) {
	active, _ := data.GetUserFromSession(in)
	return MakePostView(post, active), active
}

// Assign a priority to a post
func ScorePost(address string) int {
	stamp := strings.Split(address, "/")[1]
	result, err := strconv.Atoi(stamp)

	if err != nil {
		log.Printf("Error parsing address \"%s\"", address)
		return -1
	}

	return result
}
