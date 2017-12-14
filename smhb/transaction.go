package smhb

import (
	"container/heap"
	// "strings"
	"sync"
)

// Priority queue
type TransactionQueue struct {
	queue []*Transaction
	mut   sync.Mutex
}

// Add a post to the queue
func (this *TransactionQueue) Add(timestamp string, request REQUEST, data []byte) {
	this.mut.Lock()
	defer this.mut.Unlock()
	t := &Transaction{
		timestamp: timestamp,
		request:   request,
		data:      data,
	}

	heap.Push(this, t)
}

// Remove the post with the highest score from the queue
func (this *TransactionQueue) Remove() interface{} {
	this.mut.Lock()
	defer this.mut.Unlock()
	t := heap.Pop(this).(*Transaction)
	return t
}

// Compare two posts: higher-scored posts are closer to the top
func (this *TransactionQueue) Less(i, j int) bool {
	// TODO: compare timestamps once format is known
	this.mut.Lock()
	defer this.mut.Unlock()
	// time_i := strings.Split(this.queue[i].timestamp, "_")[0]
	// time_j := strings.Split(this.queue[j].timestamp, "_")[0]
	return true
}

/* Interface methods */

func (this *TransactionQueue) Swap(i, j int) {
	this.mut.Lock()
	defer this.mut.Unlock()
	this.queue[i], this.queue[j] = this.queue[j], this.queue[i]

	// Update indices
	this.queue[i].index = i
	this.queue[j].index = j
}

func (this *TransactionQueue) Push(x interface{}) {
	this.mut.Lock()
	defer this.mut.Unlock()
	item := x.(*Transaction)
	item.index = len(this.queue)
	this.queue = append(this.queue, item)
}

func (this *TransactionQueue) Pop() interface{} {
	this.mut.Lock()
	old := this
	index := len(old.queue)

	item := old.queue[index-1]
	this.queue = old.queue[0 : index-1]

	old.mut.Unlock()
	return item
}

func (this *TransactionQueue) Peek() interface{} {
	this.mut.Lock()
	defer this.mut.Unlock()
	item := this.queue[0]
	return item
}

func (this *TransactionQueue) Len() int {
	this.mut.Lock()
	defer this.mut.Unlock()
	return len(this.queue)
}
