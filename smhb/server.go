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
			"Request: %d; Length: %d; Target: %s",
			header.request, header.length, header.target,
		)

		// Store
		if header.length > 0 {
			data := make([]byte, header.length)
			_, err = io.ReadFull(work.connection, data)

			if err != nil {
				log.Printf("%s", err)
				continue
			}

			err = respondToStore(
				header.request, header.target, data, work.connection,
			)

			// Query
		} else {
			err = respondToQuery(
				header.request, header.target, work.connection,
			)
		}

		if err != nil {
			log.Printf("%s", err)
		}
	}

	log.Printf("worker finished execution")
}

func respondToStore(
	request REQUEST, target string, data []byte, connection net.Conn,
) error {
	var err error

	switch request {
	case USER:
		store := &userStore{}
		err = deserialize(store, data)

		if err != nil {
			return err
		}

		_, err = addUser(target, store.Password, store.Name)
		//case POST:
		//case ADD:
		//case BLOCK:
		//default:
	}

	if err != nil {
		return err
	}

	return setHeader(connection, request, 0, "")
}

func respondToQuery(
	request REQUEST, target string, connection net.Conn,
) error {
	var buffer []byte
	var err error

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
		err = errors.New("invalid request")
	}

	if err != nil {
		return err
	}

	err = setHeader(connection, request, uint16(len(buffer)), "")

	if err != nil {
		return err
	}

	_, err = connection.Write(buffer)

	return err
}
