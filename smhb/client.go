package smhb

import (
	"errors"
	"io"
	"log"
	"net"
	"strconv"
)

type Client interface {
	ServerAddress() string
	ServerPort() int
	Protocol() PROTOCOL

	GetUser(string) (User, error)
	GetPool(string) (Pool, error)
	GetPost(string) (Post, error)
	GetPostAddresses(string) ([]string, error)

	AddUser(string, string, string) (User, error)
}

func NewClient(
	serverAddress string, serverPort int, protocol PROTOCOL,
) Client {
	return client{
		serverAddress,
		serverPort,
		protocol,
	}
}

type client struct {
	serverAddress string
	serverPort    int
	protocol      PROTOCOL
}

func (this client) ServerAddress() string {
	return this.serverAddress
}

func (this client) ServerPort() int {
	return this.serverPort
}

func (this client) Protocol() PROTOCOL {
	return this.protocol
}

func (this client) query(request REQUEST, target string) ([]byte, error) {
	switch this.protocol {
	case TCP:
		bind := this.serverAddress + ":" + strconv.Itoa(this.serverPort)
		connection, err := net.Dial("tcp", bind)

		if err != nil {
			return nil, err
		}

		defer connection.Close()

		// Request

		err = setHeader(connection, request, 0, target)

		if err != nil {
			return nil, err
		}

		// Response

		header, err := getHeader(connection)

		if err != nil {
			return nil, err
		}

		log.Printf(
			"Response: %d; Length: %d; Target: %s",
			header.request, header.length, header.target,
		)

		// Validate

		if header.request == ERROR {
			return nil, errors.New("error") // tmp
		}

		if header.request != request {
			return nil, errors.New("invalid response")
		}

		if header.length == 0 {
			return nil, errors.New("data not found")
		}

		buffer := make([]byte, header.length)
		_, err = io.ReadFull(connection, buffer)

		if err != nil {
			return nil, err
		}

		return buffer, nil
	}

	return nil, nil
}

func (this client) store(request REQUEST, target string, data []byte) error {
	switch this.protocol {
	case TCP:
		bind := this.serverAddress + ":" + strconv.Itoa(this.serverPort)
		connection, err := net.Dial("tcp", bind)

		if err != nil {
			return err
		}

		defer connection.Close()

		// Request

		err = setHeader(connection, request, uint16(len(data)), target)

		if err != nil {
			return err
		}

		_, err = connection.Write(data)

		if err != nil {
			return err
		}

		// Response

		header, err := getHeader(connection)

		if err != nil {
			return err
		}

		log.Printf(
			"Response: %d; Length: %d; Target: %s",
			header.request, header.length, header.target,
		)

		// Validate

		if header.request == ERROR {
			return errors.New("error") // tmp
		}

		if header.request != request {
			return errors.New("invalid response")
		}

		return nil
	}

	return nil
}
