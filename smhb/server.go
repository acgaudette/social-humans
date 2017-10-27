package smhb

import (
	"errors"
	"io"
	"log"
	"net"
	"strconv"
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

		log.Printf(
			"[%d] Method: %d; Request: %d; Length: %d; Target: \"%s\"",
			id, header.method, header.request, header.length, header.target,
		)

		/* Response */

		switch header.method {
		case QUERY:
			err = respondToQuery(
				context, header.request, header.target, work.connection,
			)

		case STORE:
			// Read data to store
			buffer := make([]byte, header.length)
			_, err = io.ReadFull(work.connection, buffer)

			if handle(err) {
				continue CONNECTIONS
			}

			err = respondToStore(
				context, header.request, header.target, buffer, work.connection,
			)

		case EDIT:
			// Read edit data
			buffer := make([]byte, header.length)
			_, err = io.ReadFull(work.connection, buffer)

			if handle(err) {
				continue CONNECTIONS
			}

			err = respondToEdit(
				context, header.request, header.target, buffer, work.connection,
			)

		case DELETE:
			err = respondToDelete(
				context, header.request, header.target, work.connection,
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
	target string,
	connection net.Conn,
) error {
	var buffer []byte
	var err error

	// Load data by request
	switch request {
	case USER:
		buffer, err = loadUser(context, target)
	case POOL:
		buffer, err = loadPool(context, target)
	case POST:
		buffer, err = loadPost(context, target)
	case POST_ADDRESSES:
		buffer, err = serializePostAddresses(context, target)
	default:
		err = errors.New("invalid query request")
	}

	// Respond

	if err != nil {
		respondWithError(connection, err.Error())
		return err
	}

	err = setHeader(connection, QUERY, request, uint16(len(buffer)), "")

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
	target string,
	data []byte,
	connection net.Conn,
) error {
	var err error

	// Deserialize/validate incoming data
	tryRead := func(out interface{}) error {
		err = deserialize(out, data)

		if err != nil {
			respondWithError(connection, err.Error())
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

		err = addPost(context, target, store.Content, store.Author)

	default:
		err = errors.New("invalid store request")
	}

	// Respond

	if err != nil {
		respondWithError(connection, err.Error())
		return err
	}

	return setHeader(connection, STORE, request, 0, "")
}

// Edit existing data as per the client request
func respondToEdit(
	context serverContext,
	request REQUEST,
	target string,
	data []byte,
	connection net.Conn,
) error {
	var err error

	// Load and edit data by request
	switch request {
	case USER_NAME:
		loaded, err := getUser(context, target)

		if err != nil {
			respondWithError(connection, err.Error())
			return err
		}

		name := string(data)
		err = loaded.setName(context, name)

	case USER_PASSWORD:
		loaded, err := getUser(context, target)

		if err != nil {
			respondWithError(connection, err.Error())
			return err
		}

		password := string(data)
		err = loaded.updatePassword(context, password)

	case POOL_ADD:
		loaded, err := getPool(context, target)

		if err != nil {
			respondWithError(connection, err.Error())
			return err
		}

		handle := string(data)
		err = loaded.add(context, handle)

	case POOL_BLOCK:
		loaded, err := getPool(context, target)

		if err != nil {
			respondWithError(connection, err.Error())
			return err
		}

		handle := string(data)
		err = loaded.block(context, handle)

	case POST:
		loaded, err := getPost(context, target)

		if err != nil {
			respondWithError(connection, err.Error())
			return err
		}

		edit := &postEdit{}
		err = deserialize(edit, data)

		if err != nil {
			respondWithError(connection, err.Error())
			return err
		}

		err = loaded.update(context, edit.Title, edit.Content)

	default:
		err = errors.New("invalid edit request")
	}

	// Respond

	if err != nil {
		respondWithError(connection, err.Error())
		return err
	}

	return setHeader(connection, EDIT, request, 0, "")
}

// Delete data as per the client request
func respondToDelete(
	context serverContext,
	request REQUEST,
	target string,
	connection net.Conn,
) error {
	var err error

	// Delete data by request
	switch request {
	case USER:
		err = removeUser(context, target)
	case POST:
		err = removePost(context, target)
	}

	// Respond

	if err != nil {
		respondWithError(connection, err.Error())
		return err
	}

	return setHeader(connection, DELETE, request, 0, "")
}

// Send error message back to client
func respondWithError(connection net.Conn, message string) {
	err := setHeader(connection, STORE, ERROR, uint16(len(message)), "")

	if err != nil {
		log.Printf("%s", err)
	}

	_, err = connection.Write([]byte(message))

	if err != nil {
		log.Printf("%s", err)
	}
}
