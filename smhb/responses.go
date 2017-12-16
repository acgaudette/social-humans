package smhb

import (
	"bufio"
	"errors"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

// Send data to the client
func respondToQuery(
	request REQUEST,
	token Token,
	target string,
	data []byte,
	connection net.Conn,
	context ServerContext,
	access Access,
) error {
	var buffer []byte
	var err error

	// Load data by request
	switch request {
	// Generate a new token
	case TOKEN:
		loaded, err := getUser(target, context, access)

		if err != nil {
			respondWithError(connection, QUERY, ERR_NOT_FOUND, err.Error())
			return err
		}

		password := string(data)
		err = loaded.validate(password)

		if err != nil {
			respondWithError(connection, QUERY, ERR_AUTH, err.Error())
			return err
		}

		key, err := addToken(target, context, access)

		if err != nil {
			respondWithError(connection, QUERY, ERR, err.Error())
			return err
		}

		buffer = []byte(key.value)

	case USER:
		buffer, err = getRawUserInfo(target, context, access)

		if err != nil {
			respondWithError(connection, QUERY, ERR_NOT_FOUND, err.Error())
			return err
		}

	case POOL:
		if err, ok := authenticate(token, context, access); ok {
			buffer, err = getRawPool(target, context, access)

			if err != nil {
				respondWithError(connection, QUERY, ERR_NOT_FOUND, err.Error())
				return err
			}
		} else {
			respondWithError(connection, QUERY, ERR_AUTH, err.Error())
			return err
		}

	case POST_ADDRESSES:
		if err, ok := authenticate(token, context, access); ok {
			buffer, err = serializePostAddresses(target, context)

			if err != nil {
				respondWithError(connection, QUERY, ERR_NOT_FOUND, err.Error())
				return err
			}
		} else {
			respondWithError(connection, QUERY, ERR_AUTH, err.Error())
			return err
		}

	case POST:
		if err, ok := authenticate(token, context, access); ok {
			handle := strings.Split(target, "/")[0]

			// Get pool from the requester
			pool, err := getPool(handle, context, access)

			if err != nil {
				respondWithError(connection, QUERY, ERR, err.Error())
				return err
			}

			// Confirm that the requester has access to the requested
			if _, ok := pool.Users()[handle]; ok {
				buffer, err = getRawPost(target, context, access)

				if err != nil {
					respondWithError(connection, QUERY, ERR_NOT_FOUND, err.Error())
					return err
				}
			} else {
				err = errors.New("requester does not have access to requested pool")
				respondWithError(connection, QUERY, ERR_AUTH, err.Error())
				return err
			}
		} else {
			respondWithError(connection, QUERY, ERR_AUTH, err.Error())
			return err
		}

	case FEED:
		if err, ok := authenticate(token, context, access); ok {
			buffer, err = serializeFeed(target, context, access)

			if err != nil {
				respondWithError(connection, QUERY, ERR_NOT_FOUND, err.Error())
				return err
			}
		} else {
			respondWithError(connection, QUERY, ERR_AUTH, err.Error())
			return err
		}
	case INDEX:
		count, _ := countTransactions()
		buffer, err = serialize(count)
		if err != nil {
			respondWithError(connection, QUERY, ERR, err.Error())
			return err
		}
	case LOG:
		addr := connection.RemoteAddr().String()
		err := sendLog(access, context, addr)
		if err != nil {
			respondWithError(connection, QUERY, ERR, err.Error())
			return err
		}

	default:
		err = errors.New("invalid query request")
	}

	// Respond

	err = setHeader(connection, QUERY, request, uint16(len(buffer)), nil, "")

	if err != nil {
		return err
	}

	// Write serialized buffer to connection
	_, err = connection.Write(buffer)

	return err
}

// Store data sent from the client
func respondToStore(
	request REQUEST,
	token Token,
	target string,
	data []byte,
	connection net.Conn,
	context ServerContext,
	transactions *TransactionQueue,
	voteMap *sync.Map,
) error {
	timestamp := getTimestamp(context.address, context.port)

	transaction := transactions.Add(timestamp, STORE, request, target, data)
	trVote := Vote{}
	trVote.timestamp = timestamp
	trVote.finished = make(chan int)
	voteMap.Store(timestamp, &trVote)

	for _, replica := range replicas {
		go sendTransactionAction(PROPOSE, transaction, replica)
		go timeoutTransaction(&trVote, 30)
	}

	votes := <-trVote.finished
	if votes > len(replicas)/2 {
		for _, replica := range replicas {
			go sendTimestampAction(COMMIT, transaction, replica)
		}
	} else {
		log.Printf("failed to achieve quorum")
		transactions.Delete(timestamp)
		voteMap.Delete(timestamp)
		respondWithError(connection, STORE, ERR, "failed to store transaction")
	}

	// Respond
	return setHeader(connection, STORE, request, 0, nil, "")
}

// Edit existing data as per the client request
func respondToEdit(
	request REQUEST,
	token Token,
	target string,
	data []byte,
	connection net.Conn,
	context ServerContext,
	transactions *TransactionQueue,
	voteMap *sync.Map,
) error {
	timestamp := getTimestamp(context.address, context.port)

	transaction := transactions.Add(timestamp, EDIT, request, target, data)
	trVote := Vote{}
	trVote.timestamp = timestamp
	trVote.finished = make(chan int)
	voteMap.Store(timestamp, &trVote)

	for _, replica := range replicas {
		go sendTransactionAction(PROPOSE, transaction, replica)
		go timeoutTransaction(&trVote, 30)
	}

	votes := <-trVote.finished
	if votes > len(replicas)/2 {
		for _, replica := range replicas {
			go sendTimestampAction(COMMIT, transaction, replica)
		}
	} else {
		log.Printf("failed to achieve quorum")
		transactions.Delete(timestamp)
		voteMap.Delete(timestamp)
		respondWithError(connection, EDIT, ERR, "failed to store transaction")
	}

	// Respond
	return setHeader(connection, EDIT, request, 0, nil, "")
}

// Delete data as per the client request
func respondToDelete(
	request REQUEST,
	token Token,
	target string,
	connection net.Conn,
	context ServerContext,
	transactions *TransactionQueue,
	voteMap *sync.Map,
) error {
	timestamp := getTimestamp(context.address, context.port)

	transaction := transactions.Add(timestamp, DELETE, request, target, []byte{})
	trVote := Vote{}
	trVote.timestamp = timestamp
	trVote.finished = make(chan int)
	voteMap.Store(timestamp, &trVote)

	for _, replica := range replicas {
		go sendTransactionAction(PROPOSE, transaction, replica)
		go timeoutTransaction(&trVote, 30)
	}

	votes := <-trVote.finished
	if votes > len(replicas)/2 {
		for _, replica := range replicas {
			go sendTimestampAction(COMMIT, transaction, replica)
		}
	} else {
		log.Printf("failed to achieve quorum")
		transactions.Delete(timestamp)
		voteMap.Delete(timestamp)
		respondWithError(connection, DELETE, ERR, "failed to store transaction")
	}

	// Respond
	return setHeader(connection, DELETE, request, 0, nil, "")
}

// Check if data exists on the server
func respondToCheck(
	request REQUEST,
	token Token,
	target string,
	data []byte,
	connection net.Conn,
	context ServerContext,
	access Access,
) error {
	var buffer []byte
	var err error

	// Check by request
	switch request {
	case VALIDATE:
		loaded, err := getUser(target, context, access)

		if err != nil {
			respondWithError(connection, CHECK, ERR_NOT_FOUND, err.Error())
			return err
		}

		password := string(data)
		err = loaded.validate(password)

		if err != nil {
			respondWithError(connection, CHECK, ERR_AUTH, err.Error())
			return err
		}

	case TOKEN:
		_, err := getToken(target, context, access)

		if err != nil {
			respondWithError(connection, CHECK, ERR_NOT_FOUND, err.Error())
			return err
		}

	case USER:
		_, err = getRawUserInfo(target, context, access)

		if err != nil {
			respondWithError(connection, CHECK, ERR_NOT_FOUND, err.Error())
			return err
		}

	default:
		err = errors.New("invalid check request")
		respondWithError(connection, CHECK, ERR, err.Error())
		return err
	}

	// Respond

	err = setHeader(connection, CHECK, request, uint16(len(buffer)), nil, "")

	if err != nil {
		return err
	}

	// Write serialized buffer to connection
	_, err = connection.Write(buffer)

	return err
}

// Send error message back to client
func respondWithError(
	connection net.Conn,
	method METHOD,
	error REQUEST,
	message string,
) {
	err := setHeader(
		connection,
		method,
		error,
		uint16(len(message)),
		nil,
		"",
	)

	if err != nil {
		log.Printf("%s", err)
	}

	_, err = connection.Write([]byte(message))

	if err != nil {
		log.Printf("%s", err)
	}
}

func respondToPropose(
	request REQUEST,
	token Token,
	target string,
	data []byte,
	connection net.Conn,
	context ServerContext,
	access Access,
	transactions *TransactionQueue,
) error {
	var transaction *Transaction
	err := deserialize(transaction, data)
	if err != nil {
		return err
	}
	transaction.ready = make(chan bool)
	transactions.Push(transaction)

	// Wait for transaction to be first in queue
	<-transaction.ready
	err = sendTimestampAction(ACK, transaction, connection.RemoteAddr().String())
	return err
}

func respondToAck(
	request REQUEST,
	token Token,
	target string,
	data []byte,
	connection net.Conn,
	context ServerContext,
	access Access,
	transactions *TransactionQueue,
	votes *sync.Map,
) error {
	var timestamp *string
	err := deserialize(timestamp, data)
	if err != nil {
		return err
	}

	mapVal, found := votes.Load(*timestamp)
	if !found {
		return errors.New("could not find ongoing vote by timestamp")
	}
	vote := mapVal.(*Vote)

	vote.mut.Lock()
	defer vote.mut.Unlock()
	vote.votes += 1
	if vote.votes > len(replicas)/2 {
		vote.finished <- vote.votes
	}

	return nil
}

// Deserialize/validate incoming data
func tryRead(out interface{}, data []byte) error {
	err := deserialize(out, data)

	if err != nil {
		return err
	}

	return nil
}

func respondToCommit(
	token Token,
	data []byte,
	connection net.Conn,
	context ServerContext,
	access Access,
	transactions *TransactionQueue,
) error {

	var timestamp *string
	if err := tryRead(timestamp, data); err != nil {
		return err
	}

	tr := transactions.Remove()
	if tr.timestamp != *timestamp {
		return errors.New("attempted to commit transaction out of order!")
	}

	switch tr.method {
	case STORE:
		err := storeTransaction(token, connection, context, access, tr)
		if err != nil {
			return err
		}
		return logTransaction(tr, access, context)
	case EDIT:
		err := editTransaction(token, connection, context, access, tr)
		if err != nil {
			return err
		}
		return logTransaction(tr, access, context)
	case DELETE:
		err := deleteTransaction(token, connection, context, access, tr)
		if err != nil {
			return err
		}
		return logTransaction(tr, access, context)
	default:
		return errors.New("unknown METHOD for transaction")
	}
}

func storeTransaction(
	token Token,
	connection net.Conn,
	context ServerContext,
	access Access,
	tr *Transaction,
) error {
	// Store data by request
	switch tr.request {
	case USER:
		store := &userStore{}

		if err := tryRead(store, tr.data); err != nil {
			return err
		}

		_, err := addUser(tr.target, store.Password, store.Name, context, access)

		if err != nil {
			respondWithError(connection, STORE, ERR, err.Error())
			return err
		}

	case POST:
		store := &postStore{}

		if err := tryRead(store, tr.data); err != nil {
			respondWithError(connection, STORE, ERR, err.Error())
			return err
		}

		if err, ok := authenticate(token, context, access); ok {
			err = addPost(tr.target, store.Content, store.Author, context, access)

			if err != nil {
				respondWithError(connection, STORE, ERR, err.Error())
				return err
			}
		} else {
			respondWithError(connection, STORE, ERR_AUTH, err.Error())
			return err
		}

	default:
		err := errors.New("invalid store request")
		return err
	}
	return nil
}

func editTransaction(
	token Token,
	connection net.Conn,
	context ServerContext,
	access Access,
	tr *Transaction,
) error {
	// Load and edit data by request
	switch tr.request {
	case USER_NAME:
		if err, ok := authenticate(token, context, access); ok {
			loaded, err := getUser(tr.target, context, access)

			if err != nil {
				respondWithError(connection, EDIT, ERR_NOT_FOUND, err.Error())
				return err
			}

			name := string(tr.data)
			err = loaded.setName(name, context, access)

			if err != nil {
				respondWithError(connection, EDIT, ERR, err.Error())
				return err
			}
		} else {
			respondWithError(connection, EDIT, ERR_AUTH, err.Error())
			return err
		}

	case USER_PASSWORD:
		if err, ok := authenticate(token, context, access); ok {
			loaded, err := getUser(tr.target, context, access)

			if err != nil {
				respondWithError(connection, EDIT, ERR_NOT_FOUND, err.Error())
				return err
			}

			password := string(tr.data)
			err = loaded.updatePassword(password, context, access)

			if err != nil {
				respondWithError(connection, EDIT, ERR, err.Error())
				return err
			}
		} else {
			respondWithError(connection, EDIT, ERR_AUTH, err.Error())
			return err
		}

	case POOL_ADD:
		if err, ok := authenticate(token, context, access); ok {
			loaded, err := getPool(tr.target, context, access)

			if err != nil {
				respondWithError(connection, EDIT, ERR_NOT_FOUND, err.Error())
				return err
			}

			handle := string(tr.data)
			err = loaded.add(handle, context, access)

			if err != nil {
				respondWithError(connection, EDIT, ERR, err.Error())
				return err
			}
		} else {
			respondWithError(connection, EDIT, ERR_AUTH, err.Error())
			return err
		}

	case POOL_BLOCK:
		if err, ok := authenticate(token, context, access); ok {
			loaded, err := getPool(tr.target, context, access)

			if err != nil {
				respondWithError(connection, EDIT, ERR_NOT_FOUND, err.Error())
				return err
			}

			handle := string(tr.data)
			err = loaded.block(handle, context, access)

			if err != nil {
				respondWithError(connection, EDIT, ERR, err.Error())
				return err
			}
		} else {
			respondWithError(connection, EDIT, ERR_AUTH, err.Error())
			return err
		}

	case POST:
		if err, ok := authenticate(token, context, access); ok {
			loaded, err := getPost(tr.target, context, access)

			if err != nil {
				respondWithError(connection, EDIT, ERR_NOT_FOUND, err.Error())
				return err
			}

			edit := &postEdit{}
			err = deserialize(edit, tr.data)

			if err != nil {
				respondWithError(connection, EDIT, ERR, err.Error())
				return err
			}

			err = loaded.update(edit.Title, edit.Content, context, access)

			if err != nil {
				respondWithError(connection, EDIT, ERR, err.Error())
				return err
			}
		} else {
			respondWithError(connection, EDIT, ERR_AUTH, err.Error())
			return err
		}

	default:
		err := errors.New("invalid edit request")
		respondWithError(connection, EDIT, ERR, err.Error())
		return err
	}
	return nil
}

func respondToReplay(
	token Token,
	data []byte,
	connection net.Conn,
	context ServerContext,
	access Access,
) error {
	var transaction *Transaction
	err := deserialize(transaction, data)
	if err != nil {
		return err
	}

	err = logTransaction(transaction, access, context)
	if err != nil {
		return err
	}

	return nil
}

func deleteTransaction(
	token Token,
	connection net.Conn,
	context ServerContext,
	access Access,
	tr *Transaction,
) error {
	// Delete data by request
	switch tr.request {
	case USER:
		if err, ok := authenticate(token, context, access); ok {
			err = removeUser(tr.target, context, access)

			if err != nil {
				respondWithError(connection, DELETE, ERR, err.Error())
				return err
			}
		} else {
			respondWithError(connection, DELETE, ERR_AUTH, err.Error())
			return err
		}

	case POST:
		if err, ok := authenticate(token, context, access); ok {
			err = removePost(tr.target, context, access)

			if err != nil {
				respondWithError(connection, DELETE, ERR, err.Error())
				return err
			}
		} else {
			respondWithError(connection, DELETE, ERR_AUTH, err.Error())
			return err
		}
	}
	return nil
}

func sendLog(access Access, context ServerContext, destination string) error {
	file, err := os.Open("transactions/transaction.log")
	defer file.Close()
	if err != nil {
		return err
	}

	fs := bufio.NewScanner(file)
	lines := 0
	for fs.Scan() {
		tr := fs.Text()
		var transaction *Transaction
		err := access.Load(transaction, context)
		if err != nil {
			log.Printf("%s", err)
		}
		sendTransactionAction(REPLAY, transaction, destination)
	}
	return nil
}
