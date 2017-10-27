package smhb

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
)

type Client interface {
	// Configuration
	ServerAddress() string
	ServerPort() int
	Protocol() PROTOCOL

	// Queries
	GetUser(string) (User, error)
	GetPool(string) (Pool, error)
	GetPost(string) (Post, error)
	GetPostAddresses(string) ([]string, error)

	// Stores
	AddUser(string, string, string) (User, error)
	AddPost(string, string, string) error

	// Edits
	EditUserName(string, string) error
	EditUserPassword(string, string) error
	EditPoolAdd(string, string) error
	EditPoolBlock(string, string) error
	EditPost(string, string, string) error

	// Deletes
	DeleteUser(string) error
	DeletePost(string) error
}

func NewClient(
	serverAddress string, serverPort int, protocol PROTOCOL,
) Client {
	return client{
		serverAddress, serverPort, protocol,
	}
}

// Backend client
type client struct {
	serverAddress string
	serverPort    int
	protocol      PROTOCOL
}

// Interface getter methods

func (this client) ServerAddress() string {
	return this.serverAddress
}

func (this client) ServerPort() int {
	return this.serverPort
}

func (this client) Protocol() PROTOCOL {
	return this.protocol
}

// Initialize TCP connection to server
func (this client) initTCP() (net.Conn, error) {
	bind := this.serverAddress + ":" + strconv.Itoa(this.serverPort)
	connection, err := net.Dial("tcp", bind)

	if err != nil {
		return nil, err
	}

	return connection, nil
}

// Request data from the server
func (this client) query(request REQUEST, target string) ([]byte, error) {
	switch this.protocol {
	case TCP:
		connection, err := this.initTCP()

		if err != nil {
			return nil, err
		}

		defer connection.Close()

		/* Request */

		if err = setHeader(
			connection,
			QUERY,
			request,
			0,
			target,
		); err != nil {
			return nil, err
		}

		/* Response */

		header, err := getHeader(connection)

		if err != nil {
			return nil, err
		}

		/* Validate */

		err = validate(QUERY, request, header, connection)

		if err != nil {
			return nil, err
		}

		// Check for empty response
		if header.length == 0 {
			return nil, errors.New("data not returned")
		}

		/* Read data */

		buffer := make([]byte, header.length)
		_, err = io.ReadFull(connection, buffer)

		if err != nil {
			return nil, err
		}

		return buffer, nil
	}

	return nil, nil
}

// Send data to the server
func (this client) store(request REQUEST, target string, data []byte) error {
	switch this.protocol {
	case TCP:
		connection, err := this.initTCP()

		if err != nil {
			return err
		}

		defer connection.Close()

		/* Request */

		if err = setHeader(
			connection,
			STORE,
			request,
			uint16(len(data)),
			target,
		); err != nil {
			return err
		}

		// Write store buffer to connection
		_, err = connection.Write(data)

		if err != nil {
			return err
		}

		/* Response */

		header, err := getHeader(connection)

		if err != nil {
			return err
		}

		/* Validate */

		return validate(STORE, request, header, connection)
	}

	return nil
}

// Edit existing data on the server
func (this client) edit(request REQUEST, target string, data []byte) error {
	switch this.protocol {
	case TCP:
		connection, err := this.initTCP()

		if err != nil {
			return err
		}

		defer connection.Close()

		/* Request */

		if err = setHeader(
			connection,
			EDIT,
			request,
			uint16(len(data)),
			target,
		); err != nil {
			return err
		}

		// Write edit buffer to connection
		_, err = connection.Write(data)

		if err != nil {
			return err
		}

		/* Response */

		header, err := getHeader(connection)

		if err != nil {
			return err
		}

		/* Validate */

		return validate(EDIT, request, header, connection)
	}

	return nil
}

// Request data deletion from the server
func (this client) delete(request REQUEST, target string) error {
	switch this.protocol {
	case TCP:
		connection, err := this.initTCP()

		if err != nil {
			return err
		}

		defer connection.Close()

		/* Request */

		if err = setHeader(
			connection,
			DELETE,
			request,
			0,
			target,
		); err != nil {
			return err
		}

		/* Response */

		header, err := getHeader(connection)

		if err != nil {
			return err
		}

		/* Validate */

		return validate(DELETE, request, header, connection)
	}

	return nil
}

// Check for proper server response
func validate(
	method METHOD, request REQUEST, response header, connection net.Conn,
) error {
	// Check for error response
	if response.request == ERROR {
		// Read message
		buffer := make([]byte, response.length)
		_, err := io.ReadFull(connection, buffer)

		if err != nil {
			return err
		}

		message := string(buffer)
		return errors.New(message)
	}

	// Compare method
	if response.method != method {
		return errors.New(fmt.Sprintf("invalid method: %d", response.method))
	}

	// Compare request
	if response.request != request {
		return errors.New(fmt.Sprintf("invalid response: %d", response.request))
	}

	return nil
}
