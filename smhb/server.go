package smhb

import (
	"errors"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
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
) server {
	return server{
		address,
		port,
		protocol,
		poolSize,
		ServerContext{dataPath, address, port},
		NewFileAccess(),
		&TransactionQueue{},
		&sync.Map{},
	}
}

// Context for data IO
type ServerContext struct {
	dataPath string
	address  string
	port     int
}

func NewServerContext(dataPath, address string, port int) ServerContext {
	return ServerContext{dataPath, address, port}
}

// Backend server
type server struct {
	address  string
	port     int
	protocol PROTOCOL
	poolSize int
	context  ServerContext
	access   Access
	t_pq     *TransactionQueue
	votes    *sync.Map
}

type Transaction struct {
	timestamp string
	method    METHOD
	request   REQUEST
	target    string
	data      []byte
	ready     chan bool
	index     int
}

type Vote struct {
	timestamp string
	votes     int
	finished  chan int
	mut       sync.Mutex
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
		go worker(i, jobs, this.context, this.access, this.t_pq, this.votes)
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

	default:
		return errors.New("unknown server protocol")
	}

	return nil
}

type job struct {
	connection net.Conn
}

// Connection handler
func worker(id int, jobs <-chan job, context ServerContext, access Access, transactions *TransactionQueue, votes *sync.Map) {
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
				header.request,
				header.token,
				header.target,
				buffer,
				work.connection,
				context,
				access,
			)

		case STORE:
			err = respondToStore(
				header.request,
				header.token,
				header.target,
				buffer,
				work.connection,
				context,
				transactions,
				votes,
			)

		case EDIT:
			err = respondToEdit(
				header.request,
				header.token,
				header.target,
				buffer,
				work.connection,
				context,
				transactions,
				votes,
			)

		case DELETE:
			err = respondToDelete(
				header.request,
				header.token,
				header.target,
				work.connection,
				context,
				transactions,
				votes,
			)

		case CHECK:
			err = respondToCheck(
				header.request,
				header.token,
				header.target,
				buffer,
				work.connection,
				context,
				access,
			)
		case PROPOSE:
			err = respondToPropose(
				header.request,
				header.token,
				header.target,
				buffer,
				work.connection,
				context,
				access,
				transactions,
			)
		case ACK:
			err = respondToAck(
				header.request,
				header.token,
				header.target,
				buffer,
				work.connection,
				context,
				access,
				transactions,
				votes,
			)
		case COMMIT:
			err = respondToCommit(
				header.token,
				buffer,
				work.connection,
				context,
				access,
				transactions,
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

// Returns the current time via SNTP
func getNTPTime() string {
	// TODO: implement SNTP
	return ""
}

// Returns a transaction timestamp
func getTimestamp(address string, port int) string {
	return getNTPTime() + "_" + strconv.Itoa(port) + ":" + address
}
