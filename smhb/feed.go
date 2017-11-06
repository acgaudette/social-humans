package smhb

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// Internal feed representation structure
type feed struct {
	Content []string
}

// Interface getter
func (this *feed) Addresses() []string {
	return this.Content
}

// Aggregate content and create feed for a given user
func buildFeed(
	handle string, context ServerContext, access Access,
) (*feed, error) {
	// Create empty feed
	out := &feed{[]string{}}

	pool, err := getPool(handle, context, access)

	if err != nil {
		return nil, fmt.Errorf("error while accessing feed: %s", err)
	}

	q := FeedQueue{}

	// Iterate through pool and get the user posts
	for _, handle := range pool.Users() {
		addresses, err := getPostAddresses(handle, context)

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

	log.Printf("Built feed for \"%s\"", handle)

	return out, nil
}

// Create feed and serialize to buffer with lookup handle
func serializeFeed(
	handle string, context ServerContext, access Access,
) ([]byte, error) {
	out, err := buildFeed(handle, context, access)

	if err != nil {
		return nil, err
	}

	return serialize(out)
}

// Deserialize feed from buffer
func deserializeFeed(buffer []byte) (Feed, error) {
	loaded := &feed{}
	err := deserialize(loaded, buffer)

	if err != nil {
		return nil, err
	}

	return loaded, nil
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
