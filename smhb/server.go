package smhb

import (
	"errors"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type Server interface {
	// Configuration
	Address() string
	Port() int
	Protocol() PROTOCOL
	DataPath() string
	PoolSize() int

	// Operation
	ListenAndServe() error
}

func NewServer(
	address string,
	port int,
	protocol PROTOCOL,
	poolSize int,
	dataPath string,
) Server {
	return server{
		address, port, protocol, poolSize, serverContext{dataPath},
	}
}

// Context for data IO
type serverContext struct {
	dataPath string
}

// Backend server
type server struct {
	address  string
	port     int
	protocol PROTOCOL
	poolSize int
	context  serverContext
}

// Interface getter methods

func (this server) Address() string {
	return this.address
}

func (this server) Port() int {
	return this.port
}

func (this server) DataPath() string {
	return this.context.dataPath
}

func (this server) Protocol() PROTOCOL {
	return this.protocol
}

func (this server) PoolSize() int {
	return this.poolSize
}

// Handle requests and serve responses
func (this server) ListenAndServe() error {
	jobs := make(chan job, QUEUE_SIZE)

	// Spawn workers
	for i := 0; i < this.poolSize; i++ {
		go worker(this.context, i, jobs)
	}

	switch this.protocol {
	case TCP:
		bind := this.address + ":" + strconv.Itoa(this.port)
		in, err := net.Listen("tcp", bind)

		if err != nil {
			return err
		}

		defer in.Close()
		log.Printf("Listening on tcp://%s", bind)

		// Accept connections and feed them to the worker pool
		for {
			connection, err := in.Accept()

			// Handle error without halting server
			if err != nil {
				log.Printf("%s", err)
				continue
			}

			jobs <- job{connection}
		}
	}

	return nil
}

type job struct {
	connection net.Conn
}

// Connection handler
func worker(context serverContext, id int, jobs <-chan job) {
CONNECTIONS:
	// Handle a new request
	for work := range jobs {
		end := func() {
			work.connection.Close()
			log.Printf("[%d] Closed", id)
		}

		// Error handling helper
		handle := func(err error) bool {
			if err != nil {
				log.Printf("%s", err)
				end()
				return true
			}

			return false
		}

		/* Request */

		header, err := getHeader(work.connection)

		if handle(err) {
			continue CONNECTIONS
		}

		// Read data if it exists
		var buffer []byte

		if header.length > 0 {
			buffer = make([]byte, header.length)
			_, err = io.ReadFull(work.connection, buffer)

			if handle(err) {
				continue CONNECTIONS
			}
		}

		log.Printf(
			"[%d] %s %s \"%s\"; Length: %d",
			id,
			header.method.ToString(),
			header.request.ToString(),
			header.target,
			header.length,
		)

		/* Response */

		switch header.method {
		case QUERY:
			err = respondToQuery(
				context,
				header.request,
				header.token,
				header.target,
				buffer,
				work.connection,
			)

		case STORE:
			err = respondToStore(
				context,
				header.request,
				header.token,
				header.target,
				buffer,
				work.connection,
			)

		case EDIT:
			err = respondToEdit(
				context,
				header.request,
				header.token,
				header.target,
				buffer,
				work.connection,
			)

		case DELETE:
			err = respondToDelete(
				context,
				header.request,
				header.token,
				header.target,
				work.connection,
			)

		case CHECK:
			err = respondToCheck(
				context,
				header.request,
				header.token,
				header.target,
				buffer,
				work.connection,
			)
		}

		// Handle final error and close connection

		if err != nil {
			log.Printf("%s", err)
		}

		if !handle(err) {
			end()
		}
	}

	log.Printf("[%d] Execution terminated", id)
}

// Send data to the client
func respondToQuery(
	context serverContext,
	request REQUEST,
	token Token,
	target string,
	data []byte,
	connection net.Conn,
) error {
	var buffer []byte
	var err error

	// Load data by request
	switch request {
	// Generate a new token
	case TOKEN:
		loaded, err := getUser(context, target)

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

		key, err := addToken(context, target)

		if err != nil {
			respondWithError(connection, QUERY, ERR, err.Error())
			return err
		}

		buffer = []byte(key.value)

	case USER:
		buffer, err = loadUserInfo(context, target)

	case POOL:
		if err, ok := authenticate(token, target, context); ok {
			buffer, err = loadPool(context, target)
		} else {
			respondWithError(connection, QUERY, ERR_AUTH, err.Error())
			return err
		}

	case POST_ADDRESSES:
		if err, ok := authenticate(token, target, context); ok {
			buffer, err = serializePostAddresses(context, target)
		} else {
			respondWithError(connection, QUERY, ERR_AUTH, err.Error())
			return err
		}

	case POST:
		if err, ok := authenticate(token, target, context); ok {
			address := string(data)
			handle := strings.Split(address, "/")[0]

			// Get pool from the requester
			pool, err := getPool(context, target)

			if err != nil {
				respondWithError(connection, QUERY, ERR, err.Error())
				return err
			}

			// Confirm that the requester has access to the requested
			if _, ok := pool.Users()[handle]; ok {
				buffer, err = loadPost(context, address)
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
		if err, ok := authenticate(token, target, context); ok {
			buffer, err = serializeFeed(context, target)
		} else {
			respondWithError(connection, QUERY, ERR_AUTH, err.Error())
			return err
		}

	default:
		err = errors.New("invalid query request")
	}

	// Respond

	if err != nil {
		respondWithError(connection, QUERY, ERR_NOT_FOUND, err.Error())
		return err
	}

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
	context serverContext,
	request REQUEST,
	token Token,
	target string,
	data []byte,
	connection net.Conn,
) error {
	var err error

	// Deserialize/validate incoming data
	tryRead := func(out interface{}) error {
		err = deserialize(out, data)

		if err != nil {
			respondWithError(connection, STORE, ERR, err.Error())
			return err
		}

		return nil
	}

	// Store data by request
	switch request {
	case USER:
		store := &userStore{}

		if err = tryRead(store); err != nil {
			return err
		}

		_, err = addUser(context, target, store.Password, store.Name)

	case POST:
		store := &postStore{}

		if err = tryRead(store); err != nil {
			return err
		}

		if err, ok := authenticate(token, store.Author, context); ok {
			err = addPost(context, target, store.Content, store.Author)
		} else {
			respondWithError(connection, STORE, ERR_AUTH, err.Error())
			return err
		}

	default:
		err = errors.New("invalid store request")
	}

	// Respond

	if err != nil {
		respondWithError(connection, STORE, ERR, err.Error())
		return err
	}

	return setHeader(connection, STORE, request, 0, nil, "")
}

// Edit existing data as per the client request
func respondToEdit(
	context serverContext,
	request REQUEST,
	token Token,
	target string,
	data []byte,
	connection net.Conn,
) error {
	var err error

	// Load and edit data by request
	switch request {
	case USER_NAME:
		if err, ok := authenticate(token, target, context); ok {
			loaded, err := getUser(context, target)

			if err != nil {
				respondWithError(connection, EDIT, ERR_NOT_FOUND, err.Error())
				return err
			}

			name := string(data)
			err = loaded.setName(context, name)
		} else {
			respondWithError(connection, EDIT, ERR_AUTH, err.Error())
			return err
		}

	case USER_PASSWORD:
		if err, ok := authenticate(token, target, context); ok {
			loaded, err := getUser(context, target)

			if err != nil {
				respondWithError(connection, EDIT, ERR_NOT_FOUND, err.Error())
				return err
			}

			password := string(data)
			err = loaded.updatePassword(context, password)
		} else {
			respondWithError(connection, EDIT, ERR_AUTH, err.Error())
			return err
		}

	case POOL_ADD:
		if err, ok := authenticate(token, target, context); ok {
			loaded, err := getPool(context, target)

			if err != nil {
				respondWithError(connection, EDIT, ERR_NOT_FOUND, err.Error())
				return err
			}

			handle := string(data)
			err = loaded.add(context, handle)
		} else {
			respondWithError(connection, EDIT, ERR_AUTH, err.Error())
			return err
		}

	case POOL_BLOCK:
		if err, ok := authenticate(token, target, context); ok {
			loaded, err := getPool(context, target)

			if err != nil {
				respondWithError(connection, EDIT, ERR_NOT_FOUND, err.Error())
				return err
			}

			handle := string(data)
			err = loaded.block(context, handle)
		} else {
			respondWithError(connection, EDIT, ERR_AUTH, err.Error())
			return err
		}

	case POST:
		handle := strings.Split(target, "/")[0]
		if err, ok := authenticate(token, handle, context); ok {
			loaded, err := getPost(context, target)

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

			err = loaded.update(context, edit.Title, edit.Content)
		} else {
			respondWithError(connection, EDIT, ERR_AUTH, err.Error())
			return err
		}

	default:
		err = errors.New("invalid edit request")
	}

	// Respond

	if err != nil {
		respondWithError(connection, EDIT, ERR, err.Error())
		return err
	}

	return setHeader(connection, EDIT, request, 0, nil, "")
}

// Delete data as per the client request
func respondToDelete(
	context serverContext,
	request REQUEST,
	token Token,
	target string,
	connection net.Conn,
) error {
	var err error

	// Delete data by request
	switch request {
	case USER:
		if err, ok := authenticate(token, target, context); ok {
			err = removeUser(context, target)
		} else {
			respondWithError(connection, DELETE, ERR_AUTH, err.Error())
			return err
		}

	case POST:
		handle := strings.Split(target, "/")[0]
		if err, ok := authenticate(token, handle, context); ok {
			err = removePost(context, target)
		} else {
			respondWithError(connection, DELETE, ERR_AUTH, err.Error())
			return err
		}
	}

	// Respond

	if err != nil {
		respondWithError(connection, DELETE, ERR, err.Error())
		return err
	}

	return setHeader(connection, DELETE, request, 0, nil, "")
}

// Check if data exists on the server
func respondToCheck(
	context serverContext,
	request REQUEST,
	token Token,
	target string,
	data []byte,
	connection net.Conn,
) error {
	var buffer []byte
	var err error

	// Check by request
	switch request {
	case VALIDATE:
		loaded, err := getUser(context, target)

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
		_, err := getToken(context, target)

		if err != nil {
			respondWithError(connection, CHECK, ERR_NOT_FOUND, err.Error())
			return err
		}

	case USER:
		_, err = loadUserInfo(context, target)

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
