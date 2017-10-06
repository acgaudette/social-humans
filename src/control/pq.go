package control

import (
	"container/heap"
)

type post struct {
	address string
	score   int
	index   int
}

// Priority queue
type FeedQueue []*post

// Add a post to the queue
func (this *FeedQueue) Add(address string, score int) {
	post := &post{
		address: address,
		score:   score,
	}

	heap.Push(this, post)
}

// Remove the post with the highest score from the queue
func (this *FeedQueue) Remove() string {
	post := heap.Pop(this).(*post)
	return post.address
}

// Compare two posts: higher-scored posts are closer to the top
func (this FeedQueue) Less(i, j int) bool {
	return this[i].score > this[j].score
}

/* Interface methods */

func (this FeedQueue) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]

	// Update indices
	this[i].index = i
	this[j].index = j
}

func (this *FeedQueue) Push(x interface{}) {
	item := x.(*post)
	item.index = len(*this)
	*this = append(*this, item)
}

func (this *FeedQueue) Pop() interface{} {
	old := *this
	index := len(old)

	item := old[index-1]
	*this = old[0 : index-1]

	return item
}

func (this FeedQueue) Len() int {
	return len(this)
}
