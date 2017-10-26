package smhb

import (
	"errors"
	"io"
	"log"
	"net"
	"strconv"
)

type Server interface {
	Address() string
	Port() int
	Protocol() PROTOCOL
	ListenAndServe() error
}

func NewServer(address string, port int, protocol PROTOCOL) Server {
	return server{
		address,
		port,
		protocol,
	}
}

type server struct {
	address  string
	port     int
	protocol PROTOCOL
}

// Interface getter methods

func (this server) Address() string {
	return this.address
}

func (this server) Port() int {
	return this.port
}

func (this server) Protocol() PROTOCOL {
	return this.protocol
}

func (this server) ListenAndServe() error {
	jobs := make(chan job, 128)

	for i := 0; i < WORKER_COUNT; i++ {
		go worker(jobs)
	}

	switch this.protocol {
	case TCP:
		bind := this.address + ":" + strconv.Itoa(this.port)
		log.Printf("Listening on tcp://%s", bind)
		in, err := net.Listen("tcp", bind)

		if err != nil {
			return err
		}

		defer in.Close()

		for {
			connection, err := in.Accept()

			if err != nil {
				return err
			}

			jobs <- job{connection}
		}
	}

	return nil
}

type job struct {
	connection net.Conn
}

func worker(jobs <-chan job) {
	for work := range jobs {
		defer work.connection.Close()

		header, err := getHeader(work.connection)

		if err != nil {
			log.Printf("%s", err)
			continue
		}

		log.Printf(
			"[Incoming] Method: %d; Request: %d; Length: %d; Target: %s",
			header.method, header.request, header.length, header.target,
		)

		switch header.method {
		case QUERY:
			err = respondToQuery(
				header.request, header.target, work.connection,
			)

		case STORE:
			// Read data to store
			buffer := make([]byte, header.length)
			_, err = io.ReadFull(work.connection, buffer)

			if err != nil {
				log.Printf("%s", err)
				continue
			}

			err = respondToStore(
				header.request, header.target, buffer, work.connection,
			)

		case EDIT:
			// Read edit data
			buffer := make([]byte, header.length)
			_, err = io.ReadFull(work.connection, buffer)

			if err != nil {
				log.Printf("%s", err)
				continue
			}

			err = respondToEdit(
				header.request, header.target, buffer, work.connection,
			)

		case DELETE:
			err = respondToDelete(
				header.request, header.target, work.connection,
			)
		}

		if err != nil {
			log.Printf("%s", err)
		}
	}

	log.Printf("worker execution terminated")
}

func respondToQuery(
	request REQUEST, target string, connection net.Conn,
) error {
	var buffer []byte
	var err error

	// Load data by request
	switch request {
	case USER:
		buffer, err = loadUser(target)
	case POOL:
		buffer, err = loadPool(target)
	case POST:
		buffer, err = loadPost(target)
	case POST_ADDRESSES:
		buffer, err = serializePostAddresses(target)
	default:
		err = errors.New("invalid query request")
	}

	/* Response */

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

func respondToStore(
	request REQUEST, target string, data []byte, connection net.Conn,
) error {
	var err error

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

		_, err = addUser(target, store.Password, store.Name)

	case POST:
		store := &postStore{}

		if err = tryRead(store); err != nil {
			return err
		}

		err = addPost(target, store.Content, store.Author)

	default:
		err = errors.New("invalid store request")
	}

	/* Response */

	if err != nil {
		respondWithError(connection, err.Error())
		return err
	}

	return setHeader(connection, STORE, request, 0, "")
}

func respondToEdit(
	request REQUEST, target string, data []byte, connection net.Conn,
) error {
	var err error

	// Load and edit data by request
	switch request {
	case USER_NAME:
		loaded, err := getUser(target)

		if err != nil {
			respondWithError(connection, err.Error())
			return err
		}

		name := string(data)
		err = loaded.setName(name)

	case USER_PASSWORD:
		loaded, err := getUser(target)

		if err != nil {
			respondWithError(connection, err.Error())
			return err
		}

		password := string(data)
		err = loaded.updatePassword(password)

	case POOL_ADD:
		loaded, err := getPool(target)

		if err != nil {
			respondWithError(connection, err.Error())
			return err
		}

		handle := string(data)
		err = loaded.add(handle)

	case POOL_BLOCK:
		loaded, err := getPool(target)

		if err != nil {
			respondWithError(connection, err.Error())
			return err
		}

		handle := string(data)
		err = loaded.block(handle)

	case POST:
		loaded, err := getPost(target)

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

		err = loaded.update(edit.Title, edit.Content)

	default:
		err = errors.New("invalid edit request")
	}

	/* Response */

	if err != nil {
		respondWithError(connection, err.Error())
		return err
	}

	return setHeader(connection, EDIT, request, 0, "")
}

func respondToDelete(
	request REQUEST, target string, connection net.Conn,
) error {
	var err error

	// Delete data by request
	switch request {
	case USER:
		err = removeUser(target)
	case POST:
		err = removePost(target)
	}

	/* Response */

	if err != nil {
		respondWithError(connection, err.Error())
		return err
	}

	return setHeader(connection, DELETE, request, 0, "")
}

// Send error message back to client
func respondWithError(connection net.Conn, message string) {
	err := setHeader(connection, QUERY, ERROR, uint16(len(message)), "")

	if err != nil {
		log.Printf("%s", err)
	}

	_, err = connection.Write([]byte(message))

	if err != nil {
		log.Printf("%s", err)
	}
}
