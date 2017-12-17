package smhb

import (
	"container/heap"
	"fmt"
	"strings"
	"sync"
	"time"
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

func newTransaction(
	timestamp string,
	method METHOD,
	request REQUEST,
	target string,
	data []byte,
) *Transaction {
	return &Transaction{
		timestamp,
		method,
		request,
		target,
		data,
		0,
		make(chan bool, 8),
	}
}

type transactionData struct {
	Timestamp string
	Method    METHOD
	Request   REQUEST
	Target    string
	Data      []byte
}

func readTransaction(data []byte) (*Transaction, error) {
	wrapper := &transactionData{}
	err := deserialize(wrapper, data)

	if err != nil {
		return nil, fmt.Errorf("error while reading transaction: %s", err)
	}

	transaction := newTransaction(
		wrapper.Timestamp,
		wrapper.Method,
		wrapper.Request,
		wrapper.Target,
		wrapper.Data,
	)

	return transaction, nil
}

func writeTransaction(transaction *Transaction) ([]byte, error) {
	// Wrap transaction for serialization
	wrapper := transactionData{
		transaction.Timestamp,
		transaction.Method,
		transaction.Request,
		transaction.Target,
		transaction.Data,
	}

	data, err := serialize(wrapper)

	if err != nil {
		return nil, fmt.Errorf(
			"error while serializing transaction: %s",
			err.Error(),
		)
	}

	return data, nil
}

// Priority queue
type TransactionQueue struct {
	queue []*Transaction
	mut   sync.Mutex
}

// Remove the transaction with the highest score from the queue
func (this *TransactionQueue) Remove() *Transaction {
	t := heap.Pop(this).(*Transaction)
	return t
}

func (this *TransactionQueue) Peek() *Transaction {
	this.mut.Lock()
	item := this.queue[len(this.queue)-1]
	this.mut.Unlock()

	// if nil, item was deleted - repeat
	if item == nil {
		if this.Len() > 0 {
			return this.Peek()
		}
	}

	return item
}

// Deletes a transaction specified by timestamp
// returns true if specified transaction is found and deleted, false otherwise
func (this *TransactionQueue) Delete(timestamp string) bool {
	this.mut.Lock()
	defer this.mut.Unlock()

	for i := range this.queue {
		if this.queue[i] != nil && this.queue[i].Timestamp == timestamp {
			this.queue[i] = nil
			return true
		}
	}

	return false
}

// Compare two transactions: higher-scored transactions are closer to the top
func (this *TransactionQueue) Less(i, j int) bool {
	this.mut.Lock()
	defer this.mut.Unlock()

	stamp_i := strings.SplitN(this.queue[i].Timestamp, "_", 2)
	stamp_j := strings.SplitN(this.queue[j].Timestamp, "_", 2)

	time_i, _ := time.Parse(TIMESTAMP_LAYOUT, stamp_i[0])
	time_j, _ := time.Parse(TIMESTAMP_LAYOUT, stamp_j[0])

	// if timestamps are the same, compare the address/port
	// (lexicographically for now)
	if time_i == time_j {
		return stamp_i[1] < stamp_j[1]
	}

	return time_i.Before(time_j)
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

	item := x.(*Transaction)
	item.Index = len(this.queue)
	this.queue = append(this.queue, item)

	this.mut.Unlock()

	// Notify
	if this.Len() > 0 {
		this.Peek().Ready <- true
	}
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

	// Notify
	if this.Len() > 0 {
		this.Peek().Ready <- true
	}

	return item
}

func (this *TransactionQueue) Len() int {
	this.mut.Lock()
	defer this.mut.Unlock()
	return len(this.queue)
}
