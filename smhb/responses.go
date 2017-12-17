package smhb

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io/ioutil"
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
	votes *sync.Map,
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
		count, _ := countTransactions(context)
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.BigEndian, int16(count))

		if err != nil {
			log.Printf("failed to write index count")
		}

		buffer = buf.Bytes()

		if err != nil {
			respondWithError(connection, QUERY, ERR, err.Error())
			return err
		}

	case LOG:
		err := sendLog(connection, access, context, votes)

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
	timestamp, err := getTimestamp(context.address, context.port)

	if err != nil {
		respondWithError(connection, STORE, ERR, err.Error())
		return err
	}

	transaction := newTransaction(*timestamp, STORE, request, target, data)
	err = commit(transaction, transactions, voteMap)

	if err != nil {
		respondWithError(connection, STORE, ERR, err.Error())
		return err
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
	timestamp, err := getTimestamp(context.address, context.port)

	if err != nil {
		respondWithError(connection, STORE, ERR, err.Error())
		return err
	}

	transaction := newTransaction(*timestamp, EDIT, request, target, data)
	err = commit(transaction, transactions, voteMap)

	if err != nil {
		respondWithError(connection, EDIT, ERR, err.Error())
		return err
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
	timestamp, err := getTimestamp(context.address, context.port)

	if err != nil {
		respondWithError(connection, STORE, ERR, err.Error())
		return err
	}

	transaction := newTransaction(*timestamp, DELETE, request, target, []byte{})
	err = commit(transaction, transactions, voteMap)

	if err != nil {
		respondWithError(connection, DELETE, ERR, err.Error())
		return err
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
	transaction, err := readTransaction(data)

	if err != nil {
		return err
	}

	transactions.Push(transaction)

	// Wait for transaction to be first in queue
	<-transaction.Ready

	err = sendTimestamp(connection, ACK, transaction)

	if err != nil {
		return err
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
	// Read timestamp
	timestamp := string(data[:len(data)])

	tr := transactions.Remove()

	if tr.Timestamp != timestamp {
		return errors.New("attempted to commit transaction out of order!")
	}

	err := handleTransaction(
		tr,
		connection,
		access,
		context,
		token,
		false,
	)

	if err != nil {
		// Respond with failure
		if connErr := setHeader(
			connection,
			tr.Method,
			ERR,
			0,
			&token,
			"",
		); connErr != nil {
			// Handle connection error and return commit error
			log.Printf("%s", connErr.Error())
		}
	} else {
		// Respond with success
		if connErr := setHeader(
			connection,
			tr.Method,
			tr.Request,
			0,
			&token,
			"",
		); connErr != nil {
			return connErr
		}
	}

	return err
}

func storeTransaction(
	token Token,
	auth bool,
	connection net.Conn,
	context ServerContext,
	access Access,
	tr *Transaction,
) error {
	// Store data by request
	switch tr.Request {
	case USER:
		log.Printf("storing user on %s:%d", context.address, context.port)
		store := &userStore{}

		if err := tryRead(store, tr.Data); err != nil {
			return err
		}

		_, err := addUser(tr.Target, store.Password, store.Name, context, access)

		if err != nil {
			respondWithError(connection, STORE, ERR, err.Error())
			return err
		}

	case POST:
		log.Printf("storing post on %s:%d", context.address, context.port)
		store := &postStore{}

		if err := tryRead(store, tr.Data); err != nil {
			respondWithError(connection, STORE, ERR, err.Error())
			return err
		}

		if err, ok := authenticate(token, context, access); !auth || ok {
			err = addPost(tr.Target, store.Content, store.Author, tr.Timestamp, context, access)

			if err != nil {
				respondWithError(connection, STORE, ERR, err.Error())
				return err
			}
		} else {
			respondWithError(connection, STORE, ERR_AUTH, err.Error())
			return err
		}

	case TOKEN:
		store := &Token{}
		if err := tryRead(store, tr.Data); err != nil {
			return err
		}

		_, err := copyToken(tr.Target, store.value, context, access)

		if err != nil {
			respondWithError(connection, STORE, ERR, err.Error())
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
	auth bool,
	connection net.Conn,
	context ServerContext,
	access Access,
	tr *Transaction,
) error {
	// Load and edit data by request
	switch tr.Request {
	case USER_NAME:
		if err, ok := authenticate(token, context, access); !auth || ok {
			loaded, err := getUser(tr.Target, context, access)

			if err != nil {
				respondWithError(connection, EDIT, ERR_NOT_FOUND, err.Error())
				return err
			}

			name := string(tr.Data)
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
		if err, ok := authenticate(token, context, access); !auth || ok {
			loaded, err := getUser(tr.Target, context, access)

			if err != nil {
				respondWithError(connection, EDIT, ERR_NOT_FOUND, err.Error())
				return err
			}

			password := string(tr.Data)
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
		if err, ok := authenticate(token, context, access); !auth || ok {
			loaded, err := getPool(tr.Target, context, access)

			if err != nil {
				respondWithError(connection, EDIT, ERR_NOT_FOUND, err.Error())
				return err
			}

			handle := string(tr.Data)
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
		if err, ok := authenticate(token, context, access); !auth || ok {
			loaded, err := getPool(tr.Target, context, access)

			if err != nil {
				respondWithError(connection, EDIT, ERR_NOT_FOUND, err.Error())
				return err
			}

			handle := string(tr.Data)
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
		if err, ok := authenticate(token, context, access); !auth || ok {
			loaded, err := getPost(tr.Target, context, access)

			if err != nil {
				respondWithError(connection, EDIT, ERR_NOT_FOUND, err.Error())
				return err
			}

			edit := &postEdit{}
			err = deserialize(edit, tr.Data)

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
	transaction, err := readTransaction(data)

	if err != nil {
		return err
	}

	return logTransaction(transaction, access, context)
}

func deleteTransaction(
	token Token,
	auth bool,
	connection net.Conn,
	context ServerContext,
	access Access,
	tr *Transaction,
) error {
	// Delete data by request
	switch tr.Request {
	case USER:
		if err, ok := authenticate(token, context, access); !auth || ok {
			err = removeUser(tr.Target, context, access)

			if err != nil {
				respondWithError(connection, DELETE, ERR, err.Error())
				return err
			}
		} else {
			respondWithError(connection, DELETE, ERR_AUTH, err.Error())
			return err
		}

	case POST:
		if err, ok := authenticate(token, context, access); !auth || ok {
			err = removePost(tr.Target, context, access)

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

func sendLog(
	connection net.Conn, access Access, context ServerContext, votes *sync.Map,
) error {
	file, err := os.Open(context.dataPath + "/transactions.log")

	if err != nil {
		return err
	}

	defer file.Close()

	fs := bufio.NewScanner(file)

	for fs.Scan() {
		tr := fs.Text() + ".trans"
		data, err := ioutil.ReadFile(context.dataPath + "/log/" + tr)

		if err != nil {
			log.Printf("sendLog: could not find transaction %s", tr)
			continue
		}

		transaction, err := readTransaction(data)

		if err != nil {
			log.Printf("error reading transaction: %s", err)
			continue
		}

		err = sendTransaction(connection, REPLAY, transaction)

		if err != nil {
			log.Printf("error sending transaction: %s", err)
		}
	}

	return nil
}
