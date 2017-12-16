package smhb

import (
	"net"
	"time"
)

func connect(destination string) (net.Conn, error) {
	connection, err := net.DialTimeout("tcp", destination, time.Second*20)

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
) error {
	connection, err := connect(destination)

	if err != nil {
		return ConnectionError{err}
	}

	defer connection.Close()

	// No token checking for replication processes (RIP)
	token := Token{}

	data, err := serialize(transaction)
	if err != nil {
		return err
	}

	if err = setHeader(
		connection,
		PROPOSE,
		transaction.request,
		uint16(len(data)),
		&token,
		transaction.target,
	); err != nil {
		return ConnectionError{err}
	}

	// Write store buffer to connection
	_, err = connection.Write(data)

	if err != nil {
		return ConnectionError{err}
	}

	return nil
}

// Sends a timestamp with the specified method in header
func sendTimestampAction(
	method METHOD,
	transaction *Transaction,
	destination string,
) error {
	connection, err := connect(destination)

	if err != nil {
		return ConnectionError{err}
	}

	defer connection.Close()

	// No token checking for replication processes (RIP)
	token := Token{}

	data, err := serialize(transaction.timestamp)
	if err != nil {
		return err
	}

	if err = setHeader(
		connection,
		PROPOSE,
		transaction.request,
		uint16(len(data)),
		&token,
		transaction.target,
	); err != nil {
		return ConnectionError{err}
	}

	// Write store buffer to connection
	_, err = connection.Write(data)

	if err != nil {
		return ConnectionError{err}
	}

	return nil
}
