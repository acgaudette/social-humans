package smhb

import (
	"container/heap"
)

type postptr struct {
	address string
	score   int
	index   int
}

// Priority queue
type FeedQueue []*postptr

// Add a post to the queue
func (this *FeedQueue) Add(address string, score int) {
	post := &postptr{
		address: address,
		score:   score,
	}

	heap.Push(this, post)
}

// Remove the post with the highest score from the queue
func (this *FeedQueue) Remove() string {
	post := heap.Pop(this).(*postptr)
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
	item := x.(*postptr)
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
