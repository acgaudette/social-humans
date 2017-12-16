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
	votes *sync.Map,
) {
	connection, err := connect(destination)

	if err != nil {
		log.Printf("%s", ConnectionError{err}.Error())
		return
	}

	defer connection.Close()

	// No token checking for replication processes
	token := Token{}

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

	/* Respond to ACK */

	header, err := getHeader(connection)

	if err != nil {
		log.Printf("%s", err.Error())
		return
	}

	// Read timestamp

	data = make([]byte, header.length)
	_, err = connection.Read(data[:])

	if err != nil {
		log.Printf("%s", err.Error())
		return
	}

	timestamp := string(data[:header.length])

	mapVal, found := votes.Load(timestamp)

	if !found {
		log.Printf("Could not find ongoing vote by timestamp")
		return
	}

	vote := mapVal.(*Vote)

	vote.mut.Lock()
	defer vote.mut.Unlock()

	vote.votes += 1
	if vote.votes > len(replicas)/2 {
		vote.finished <- vote.votes
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

	err = sendTimestamp(connection, method, transaction)

	if err != nil {
		log.Printf("%s", err.Error())
	}
}

func sendTimestamp(
	connection net.Conn,
	method METHOD,
	transaction *Transaction,
) error {
	// No token checking for replication processes
	token := Token{}

	if err := setHeader(
		connection,
		method,
		transaction.Request,
		uint16(len(transaction.Timestamp)),
		&token,
		transaction.Target,
	); err != nil {
		return err
	}

	// Write store buffer to connection
	_, err := connection.Write([]byte(transaction.Timestamp))

	if err != nil {
		return ConnectionError{err}
	}

	return nil
}
