package smhb

import (
	"errors"
	"io"
	"log"
	"net"
	"strconv"
	"fmt"
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
	AddPost(string, string, string) error

	EditPoolAdd(string, string) error
	EditPoolBlock(string, string) error

	DeleteUser(string) error
	DeletePost(string) error
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

		err = setHeader(connection, QUERY, request, 0, target)

		if err != nil {
			return nil, err
		}

		// Response

		header, err := getHeader(connection)

		if err != nil {
			return nil, err
		}

		log.Printf(
			"Response: %d/%d; Length: %d; Target: %s",
			header.method, header.request, header.length, header.target,
		)

		// Validate

		err = validate(QUERY, request, header, connection)

		if err != nil {
			return nil, err
		}

		// Check for empty response
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

		err = setHeader(connection, STORE, request, uint16(len(data)), target)

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
			"Response: %d/%d; Length: %d; Target: %s",
			header.method, header.request, header.length, header.target,
		)

		// Validate

		return validate(STORE, request, header, connection)
	}

	return nil
}

func (this client) edit(request REQUEST, target string, data []byte) error {
	switch this.protocol {
	case TCP:
		bind := this.serverAddress + ":" + strconv.Itoa(this.serverPort)
		connection, err := net.Dial("tcp", bind)

		if err != nil {
			return err
		}

		defer connection.Close()

		// Request

		err = setHeader(connection, EDIT, request, uint16(len(data)), target)

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
			"Response: %d/%d; Length: %d; Target: %s",
			header.method, header.request, header.length, header.target,
		)

		// Validate

		return validate(EDIT, request, header, connection)
	}

	return nil
}

func (this client) delete(request REQUEST, target string) error {
	switch this.protocol {
	case TCP:
		bind := this.serverAddress + ":" + strconv.Itoa(this.serverPort)
		connection, err := net.Dial("tcp", bind)

		if err != nil {
			return err
		}

		defer connection.Close()

		// Request

		err = setHeader(connection, DELETE, request, 0, target)

		if err != nil {
			return err
		}

		// Response

		header, err := getHeader(connection)

		if err != nil {
			return err
		}

		log.Printf(
			"Response: %d/%d; Length: %d; Target: %s",
			header.method, header.request, header.length, header.target,
		)

		// Validate

		return validate(DELETE, request, header, connection)
	}

	return nil
}

func validate(
	method METHOD, request REQUEST, response header, connection net.Conn,
) error {
	// Check for error response
	if response.request == ERROR {
		buffer := make([]byte, response.length)
		_, err := io.ReadFull(connection, buffer)

		if err != nil {
			return err
		}

		message := string(buffer)
		return errors.New(message)
	}

	// Check for response mismatch

	if response.method != method {
		return errors.New(fmt.Sprintf("invalid method: %d", response.method))
	}

	if response.request != request {
		return errors.New(fmt.Sprintf("invalid response: %d", response.request))
	}

	return nil
}
