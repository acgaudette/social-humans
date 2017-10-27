package smhb

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type feed struct {
	Content []string
}

func (this *feed) Addresses() []string {
	return this.Content
}

func buildFeed(context serverContext, handle string) (*feed, error) {
	// Create empty feed
	out := &feed{[]string{}}

	pool, err := getPool(context, handle)

	if err != nil {
		return nil, fmt.Errorf("error while accessing feed: %s", err)
	}

	q := FeedQueue{}

	// Iterate through pool and get the user posts
	for _, handle := range pool.Users() {
		addresses, err := getPostAddresses(context, handle)

		if err != nil {
			log.Printf("Error getting posts from \"%s\": %s", handle, err)
			continue
		}

		// Iterate through posts and push to the priority queue
		for _, post := range addresses {
			score, err := scorePost(post)

			if err != nil {
				log.Printf("Error parsing address \"%s\": %s", post, err)
				continue
			}

			q.Add(post, score)
		}
	}

	// Convert queue into feed view
	for q.Len() > 0 {
		out.Content = append(out.Content, q.Remove())
	}

	return out, nil
}

// Assign a priority to a post
func scorePost(address string) (int, error) {
	stamp := strings.Split(address, "/")[1]
	result, err := strconv.Atoi(stamp)

	if err != nil {
		return -1, err
	}

	return result, err
}
