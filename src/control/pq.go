package control

import (
	"container/heap"
)

type Item struct {
	value    string
	priority int
	index    int
}

type PQueue []*Item

// Add an item to the queue
func (this *PQueue) Add(value string, priority int) {
	item := &Item{value: value, priority: priority}
	heap.Push(this, item)
}

// Remove the highest-priority item from the queue
func (this *PQueue) Remove() string {
	item := heap.Pop(this).(*Item)
	return item.value
}

func (this PQueue) Less(i, j int) bool {
	return this[i].priority > this[j].priority
}

func (this PQueue) Swap(i, j int) {
	// Swap
	this[i], this[j] = this[j], this[i]

	// Update indices
	this[i].index = i
	this[j].index = j
}

func (this *PQueue) Push(x interface{}) {
	item := x.(*Item)
	item.index = len(*this)
	*this = append(*this, item)
}

func (this *PQueue) Pop() interface{} {
	old := *this
	index := len(old)

	item := old[index-1]
	*this = old[0 : index-1]

	return item
}

func (this PQueue) Len() int {
	return len(this)
}
