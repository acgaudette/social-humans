package smhb

import (
	"fmt"
)

// Representation of idempotent, timestamped write updates
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

// Transaction data wrapper for storage
type transactionData struct {
	Timestamp string
	Method    METHOD
	Request   REQUEST
	Target    string
	Data      []byte
}

// Read transaction from buffer
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

// Write transaction to buffer
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
