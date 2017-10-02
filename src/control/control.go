package control

import (
	"../data"
	"../front"
	"log"
	"net/http"
	"strconv"
)

/*
	Model to view functions never return nil, so that something is always
	rendered
*/

// Build a views map from the active user and a generic view
func MakeViews(view interface{}, activeUser *data.User) *front.Views {
	views := make(front.Views)

	// Content (main) view
	if view != nil {
		views["content"] = view
	}

	// Active user account view
	if activeUser != nil {
		views["active"] = MakeActiveView(activeUser.Handle)
	}

	return &views
}

// MakeViews, but automatically load the active user
func GetUserAndMakeViews(view interface{}, in *http.Request) *front.Views {
	account, _ := data.GetUserFromSession(in) // Nil check done in MakeViews
	return MakeViews(view, account)
}

// Build an ActiveView
func MakeActiveView(handle string) *front.ActiveView {
	return &front.ActiveView{
		Handle: handle,
	}
}

// Build a UserView from a user model
func MakeUserView(
	user *data.User, status string, account *data.User,
) *front.UserView {
	handle := user.Handle

	// Always display something to the frontend
	if handle == "" {
		handle = "Username Invalid"
	}

	name := user.Name

	// Always display something to the frontend
	if name == "" {
		name = "Name Invalid"
	}

	// Compare the active user to the input user
	active := false

	if account != nil && account.Handle == user.Handle {
		active = true
	}

	return &front.UserView{
		Handle:       handle,
		Name:         name,
		Status:       status,
		IsActiveUser: active,
	}
}

// Load the active user and build a UserView
func GetUserAndMakeUserView(
	user *data.User, status string, in *http.Request,
) (*front.UserView, *data.User) {
	account, _ := data.GetUserFromSession(in)
	return MakeUserView(user, status, account), account
}

// Build a PoolView
func MakePoolView(handle string, status string) (*front.PoolView, error) {
	pool, err := data.LoadPool(handle)

	if err != nil {
		// Return empty pool view if pool is not found
		empty := &front.PoolView{
			Handles: []string{},
			Status:  "Error: access failure",
		}

		return empty, err
	}

	if len(pool.Users) <= 1 {
		// Override the empty pool message with the input status message
		if status == "" {
			status = "Your pool is empty!"
		}

		// Return empty pool view
		empty := &front.PoolView{
			Handles: []string{},
			Status:  status,
		}

		return empty, nil
	}

	result := &front.PoolView{
		Handles: []string{},
		Status:  status,
	}

	// Build handles slice from pool users
	for _, value := range pool.Users {
		if value == handle {
			continue
		}

		result.Handles = append(result.Handles, value)
	}

	return result, nil
}

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

	if len(pool.Users) <= 1 {
		// Return empty feed view if pool is empty
		feed.Status = "Your pool is empty!"
		return feed, nil
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
	for _, item := range q {
		address := item.value
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

		view := MakePostView(post)
		feed.Posts = append(feed.Posts, view)
	}

	if len(feed.Posts) == 0 {
		feed.Status = "Nothing to see here..."
	}

	return feed, nil
}

// Build a PostView from a post model
func MakePostView(post *data.Post) *front.PostView {
	return &front.PostView{
		Title:   post.Title,
		Content: post.Content,
		Author:  post.Author,
	}
}

// Assign a priority to a post
func ScorePost(address string) int {
	result, err := strconv.Atoi(address)

	if err != nil {
		return -1
	}

	return result
}
