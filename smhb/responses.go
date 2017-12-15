package smhb

import (
	"errors"
	"log"
	"net"
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
	access Access,
	transactions *TransactionQueue,
	voteMap *sync.Map,
) error {
	var err error

	timestamp := getTimestamp(context.address, context.port)

	// Deserialize/validate incoming data
	tryRead := func(out interface{}) error {
		err = deserialize(out, data)

		if err != nil {
			respondWithError(connection, STORE, ERR, err.Error())
			return err
		}

		return nil
	}

	transactions.Add(timestamp, request, target, data)
	trVote := Vote{}
	trVote.timestamp = timestamp
	trVote.finished = make(chan int)
	voteMap.Store(timestamp, trVote)

	for _, replica := range replicas {
		go proposeTransaction(request, target, data, timestamp, replica)
	}

	votes := <-trVote.finished
	if votes >= len(replicas)+1 {
		go commitTransaction(timestamp)
	}

	// Store data by request
	switch request {
	case USER:
		store := &userStore{}

		if err = tryRead(store); err != nil {
			return err
		}

		_, err = addUser(target, store.Password, store.Name, context, access)

		if err != nil {
			respondWithError(connection, STORE, ERR, err.Error())
			return err
		}

	case POST:
		store := &postStore{}

		if err = tryRead(store); err != nil {
			return err
		}

		if err, ok := authenticate(token, context, access); ok {
			err = addPost(target, store.Content, store.Author, context, access)

			if err != nil {
				respondWithError(connection, STORE, ERR, err.Error())
				return err
			}
		} else {
			respondWithError(connection, STORE, ERR_AUTH, err.Error())
			return err
		}

	default:
		err = errors.New("invalid store request")
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
	access Access,
) error {
	// Load and edit data by request
	switch request {
	case USER_NAME:
		if err, ok := authenticate(token, context, access); ok {
			loaded, err := getUser(target, context, access)

			if err != nil {
				respondWithError(connection, EDIT, ERR_NOT_FOUND, err.Error())
				return err
			}

			name := string(data)
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
			loaded, err := getUser(target, context, access)

			if err != nil {
				respondWithError(connection, EDIT, ERR_NOT_FOUND, err.Error())
				return err
			}

			password := string(data)
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
			loaded, err := getPool(target, context, access)

			if err != nil {
				respondWithError(connection, EDIT, ERR_NOT_FOUND, err.Error())
				return err
			}

			handle := string(data)
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
			loaded, err := getPool(target, context, access)

			if err != nil {
				respondWithError(connection, EDIT, ERR_NOT_FOUND, err.Error())
				return err
			}

			handle := string(data)
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
			loaded, err := getPost(target, context, access)

			if err != nil {
				respondWithError(connection, EDIT, ERR_NOT_FOUND, err.Error())
				return err
			}

			edit := &postEdit{}
			err = deserialize(edit, data)

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
	access Access,
) error {
	// Delete data by request
	switch request {
	case USER:
		if err, ok := authenticate(token, context, access); ok {
			err = removeUser(target, context, access)

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
			err = removePost(target, context, access)

			if err != nil {
				respondWithError(connection, DELETE, ERR, err.Error())
				return err
			}
		} else {
			respondWithError(connection, DELETE, ERR_AUTH, err.Error())
			return err
		}
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
	voteMap *sync.Map,
) error {
	// // Deserialize/validate incoming data
	// tryRead := func(out interface{}) error {
	// 	err = deserialize(out, data)

	// 	if err != nil {
	// 		respondWithError(connection, STORE, ERR, err.Error())
	// 		return err
	// 	}

	// 	return nil
	// }

	// transactions.Add(timestamp, request, target, data)
	// if transactions.Peek().timestamp == timestamp {
	// 	err = ackTransaction()
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
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
	voteMap *sync.Map,
) error {
	return nil
}

func respondToCommit(
	request REQUEST,
	token Token,
	target string,
	data []byte,
	connection net.Conn,
	context ServerContext,
	access Access,
	transactions *TransactionQueue,
	voteMap *sync.Map,
) error {
	return nil
}
