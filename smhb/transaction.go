package smhb

import (
	"container/heap"
	"fmt"
	"sync"
)

type Transaction struct {
	Timestamp string
	Method    METHOD
	Request   REQUEST
	Target    string
	Data      []byte
	Index     int
	Ready     chan bool
}

type transactionData struct {
	Timestamp string
	Method    METHOD
	Request   REQUEST
	Target    string
	Data      []byte
	Index     int
}

// Priority queue
type TransactionQueue struct {
	queue []*Transaction
	mut   sync.Mutex
}

// Add a transaction to the queue
func (this *TransactionQueue) Add(
	timestamp string, method METHOD, request REQUEST, target string, data []byte,
) *Transaction {
	t := &Transaction{
		Timestamp: timestamp,
		Method:    method,
		Request:   request,
		Target:    target,
		Data:      data,
		Ready:     make(chan bool),
	}

	heap.Push(this, t)
	return t
}

// Remove the transaction with the highest score from the queue
func (this *TransactionQueue) Remove() *Transaction {
	this.mut.Lock()
	defer this.mut.Unlock()
	t := heap.Pop(this).(*Transaction)
	return t
}

/*
Deletes a transaction specified by timestamp
returns true if specified transaction is found and deleted, false otherwise
*/
func (this *TransactionQueue) Delete(timestamp string) bool {
	this.mut.Lock()
	defer this.mut.Unlock()
	for i := range this.queue {
		if this.queue[i].Timestamp == timestamp {
			this.queue[i] = nil
			return true
		}
	}
	return false
}

// Compare two transactions: higher-scored transactions are closer to the top
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
	this.queue[i].Index = i
	this.queue[j].Index = j
}

func (this *TransactionQueue) Push(x interface{}) {
	this.mut.Lock()
	defer this.mut.Unlock()
	item := x.(*Transaction)
	item.Index = len(this.queue)
	this.queue = append(this.queue, item)
}

func (this *TransactionQueue) Pop() interface{} {
	this.mut.Lock()
	old := this
	index := len(old.queue)

	item := old.queue[index-1]
	this.queue = old.queue[0 : index-1]
	old.mut.Unlock()
	// if nil, item was deleted - repeat
	if item == nil {
		return this.Pop()
	}
	return item
}

func (this *TransactionQueue) Peek() interface{} {
	this.mut.Lock()
	item := this.queue[len(this.queue)-1]
	this.mut.Unlock()
	// if nil, item was deleted - repeat
	if item == nil {
		return this.Peek()
	}
	return item
}

func (this *TransactionQueue) Len() int {
	this.mut.Lock()
	defer this.mut.Unlock()
	return len(this.queue)
}
