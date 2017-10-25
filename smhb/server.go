package smhb

import (
	"encoding/binary"
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
	jobs := make(chan job, 128) //

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

func NewServer(address string, port int, protocol PROTOCOL) Server {
	return server{
		address,
		port,
		protocol,
	}
}

type job struct {
	connection net.Conn
}

func worker(jobs <-chan job) {
	for work := range jobs {
		defer work.connection.Close()

		var request, length uint16

		err := binary.Read(work.connection, binary.LittleEndian, &request)

		if err != nil {
			log.Printf("%s", err)
			continue
		}

		err = binary.Read(work.connection, binary.LittleEndian, &length)

		if err != nil {
			log.Printf("%s", err)
			continue
		}

		log.Printf("Request: %d; Length: %d", request, length)
	}
}
