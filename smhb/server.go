package smhb

import (
	"errors"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
	"fmt"
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

		// Check if catchup is needed
		go this.checkLog()

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

		case REPLAY:
			err = respondToReplay(
				header.token,
				buffer,
				work.connection,
				context,
				access,
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

// Returns a transaction timestamp
func getTimestamp(address string, port int) (*string, error) {
	result, err := getNTPTime()

	if err != nil {
		return nil, err
	}

	stamp := (*result).Format(time.RFC3339Nano) +
		"_" + strconv.Itoa(port) +
		":" + address

	return &stamp, nil
}

type maxCount struct {
	max  int
	addr string
	mut  *sync.Mutex
}

func (this server) checkLog() {
	count, _ := countTransactions()
	m := maxCount{}
	for _, replica := range replicas {
		go queryMaxIndex(&m, replica)
	}
	if m.max > count {
		go requestLog(m.addr)
	}
}

func queryMaxIndex(m *maxCount, destination string) {
	connection, err := connect(destination)

	if err != nil {
		return
	}

	if err = setHeader(
		connection,
		QUERY,
		INDEX,
		uint16(0),
		nil,
		"",
	); err != nil {
		return
	}

	header, err := getHeader(connection)
	if err != nil {
		return
	}

	if header.length == 0 {
		return
	}

	buffer := make([]byte, header.length)
	_, err = io.ReadFull(connection, buffer)

	if err != nil {
		return
	}
	var count int
	err = deserialize(count, buffer)
	if err != nil {
		return
	}

	m.mut.Lock()
	defer m.mut.Unlock()
	if count > m.max {
		m.max = count
		m.addr = destination
	}
}

func requestLog(destination string) {
	connection, err := connect(destination)

	if err != nil {
		return
	}

	if err = setHeader(
		connection,
		QUERY,
		LOG,
		uint16(0),
		nil,
		"",
	); err != nil {
		return
	}
}

func commit(
	transaction *Transaction, transactions *TransactionQueue, votes *sync.Map,
) error {
	vote := Vote{
		timestamp: transaction.timestamp,
		finished:  make(chan int),
	}

	votes.Store(transaction.timestamp, &vote)

	// Propose transaction
	for _, replica := range replicas {
		go sendTransactionAction(PROPOSE, transaction, replica)
		go timeoutTransaction(&vote, TRANSACTION_TIMEOUT)
	}

	// Attempt to acquire quorum
	count := <-vote.finished

	if count > len(replicas)/2 {
		// Success; commit
		for _, replica := range replicas {
			go sendTimestampAction(COMMIT, transaction, replica)
		}

		return nil
	}

	// Quorum not achieved; abort
	transactions.Delete(transaction.timestamp)
	votes.Delete(transaction.timestamp)
	return fmt.Errorf("failed to achieve quorum (count is %d)", count)
}
