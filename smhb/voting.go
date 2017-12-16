package smhb

import (
	"net"
	"time"
	"log"
)

func connect(destination string) (net.Conn, error) {
	connection, err := net.DialTimeout(
		"tcp",
		destination,
		RM_TIMEOUT * time.Second,
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

	data, err := serialize(transaction)

	if err != nil {
		log.Printf("%s", err.Error())
		return
	}

	if err = setHeader(
		connection,
		method,
		transaction.request,
		uint16(len(data)),
		&token,
		transaction.target,
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

	data, err := serialize(transaction.timestamp)

	if err != nil {
		log.Printf("%s", err.Error())
		return
	}

	if err = setHeader(
		connection,
		method,
		transaction.request,
		uint16(len(data)),
		&token,
		transaction.target,
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
