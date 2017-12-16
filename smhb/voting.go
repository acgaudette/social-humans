package smhb

import (
	"log"
	"net"
	"sync"
	"time"
)

type Vote struct {
	timestamp string
	votes     int
	finished  chan int
	mut       sync.Mutex
}

func connect(destination string) (net.Conn, error) {
	connection, err := net.DialTimeout(
		"tcp",
		destination,
		RM_TIMEOUT*time.Second,
	)

	if err != nil {
		return nil, err
	}

	return connection, nil
}

// Times out a vote after "timeout" seconds
func timeoutTransaction(vote *Vote, timeout int) {
	time.Sleep(time.Second * time.Duration(timeout))
	vote.finished <- vote.votes
}

// Sends an entire transaction with the specified method in header
func sendTransactionAction(
	method METHOD,
	transaction *Transaction,
	destination string,
) {
	connection, err := connect(destination)

	if err != nil {
		log.Printf("%s", ConnectionError{err}.Error())
		return
	}

	defer connection.Close()

	// No token checking for replication processes (RIP)
	token := Token{}

	// Wrap transaction for serialization
	wrapper := transactionData{
		transaction.Timestamp,
		transaction.Method,
		transaction.Request,
		transaction.Target,
		transaction.Data,
		transaction.Index,
	}

	data, err := serialize(wrapper)

	if err != nil {
		log.Printf("error while serializing transaction: %s", err.Error())
		return
	}

	if err = setHeader(
		connection,
		method,
		transaction.Request,
		uint16(len(data)),
		&token,
		transaction.Target,
	); err != nil {
		log.Printf("%s", err.Error())
		return
	}

	// Write store buffer to connection
	_, err = connection.Write(data)

	if err != nil {
		log.Printf("%s", ConnectionError{err}.Error())
		return
	}
}

// Sends a timestamp with the specified method in header
func sendTimestampAction(
	method METHOD,
	transaction *Transaction,
	destination string,
) {
	connection, err := connect(destination)

	if err != nil {
		log.Printf("%s", ConnectionError{err}.Error())
		return
	}

	defer connection.Close()

	// No token checking for replication processes (RIP)
	token := Token{}

	data, err := serialize(transaction.Timestamp)

	if err != nil {
		log.Printf("error while serializing timestamp: %s", err.Error())
		return
	}

	if err = setHeader(
		connection,
		method,
		transaction.Request,
		uint16(len(data)),
		&token,
		transaction.Target,
	); err != nil {
		log.Printf("%s", err.Error())
		return
	}

	// Write store buffer to connection
	_, err = connection.Write(data)

	if err != nil {
		log.Printf("%s", ConnectionError{err}.Error())
		return
	}
}
