package smhb

import (
	"io/ioutil"
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

		// Request

		header, err := getHeader(work.connection)

		if err != nil {
			log.Printf("%s", err)
			continue
		}

		log.Printf("Request: %d; Length: %d", header.request, header.length)

		switch header.request {
		case USER:
			handle := "acg"
			buffer, err := ioutil.ReadFile(prefix(handle + ".user"))

			if err != nil {
				log.Printf("%s", err)
				continue
			}

			log.Printf("length: %d", len(buffer))

			// Response

			err = setHeader(work.connection, header.request, uint16(len(buffer)))

			if err != nil {
				log.Printf("%s", err)
				continue
			}

			work.connection.Write(buffer)
		}
	}

	log.Printf("worker finished execution")
}
