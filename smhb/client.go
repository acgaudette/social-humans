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

	// Validation
	Validate(string, string) (bool, error)

	// Queries
	GetUser(string) (User, error)
	GetPool(string) (Pool, error)
	GetPost(string) (Post, error)
	GetPostAddresses(string) ([]string, error)
	GetFeed(string) (Feed, error)

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
func (this client) query(
	request REQUEST, target string, data []byte, token *Token,
) ([]byte, error) {
	switch this.protocol {
	case TCP:
		connection, err := this.initTCP()

		if err != nil {
			return nil, ConnectionError{err}
		}

		defer connection.Close()

		/* Request */

		length := uint16(0)

		if data != nil {
			length = uint16(len(data))
		}

		if err = setHeader(
			connection,
			QUERY,
			request,
			length,
			token,
			target,
		); err != nil {
			return nil, ConnectionError{err}
		}

		// Write edit buffer to connection
		if data != nil {
			_, err = connection.Write(data)

			if err != nil {
				return nil, ConnectionError{err}
			}
		}

		/* Response */

		header, err := getHeader(connection)

		if err != nil {
			return nil, ConnectionError{err}
		}

		/* Validate */

		smhbErr := validate(QUERY, request, header, connection)

		if smhbErr != nil {
			return nil, smhbErr
		}

		// Check for empty response
		if header.length == 0 {
			return nil, ConnectionError{errors.New("data not returned")}
		}

		/* Read data */

		buffer := make([]byte, header.length)
		_, err = io.ReadFull(connection, buffer)

		if err != nil {
			return nil, ConnectionError{err}
		}

		return buffer, nil
	}

	return nil, nil
}

// Send data to the server
func (this client) store(
	request REQUEST, target string, data []byte, token *Token,
) error {
	switch this.protocol {
	case TCP:
		connection, err := this.initTCP()

		if err != nil {
			return ConnectionError{err}
		}

		defer connection.Close()

		/* Request */

		if err = setHeader(
			connection,
			STORE,
			request,
			uint16(len(data)),
			token,
			target,
		); err != nil {
			return ConnectionError{err}
		}

		// Write store buffer to connection
		_, err = connection.Write(data)

		if err != nil {
			return ConnectionError{err}
		}

		/* Response */

		header, err := getHeader(connection)

		if err != nil {
			return ConnectionError{err}
		}

		/* Validate */

		return validate(STORE, request, header, connection)
	}

	return nil
}

// Edit existing data on the server
func (this client) edit(
	request REQUEST, target string, data []byte, token *Token,
) error {
	switch this.protocol {
	case TCP:
		connection, err := this.initTCP()

		if err != nil {
			return ConnectionError{err}
		}

		defer connection.Close()

		/* Request */

		if err = setHeader(
			connection,
			EDIT,
			request,
			uint16(len(data)),
			token,
			target,
		); err != nil {
			return ConnectionError{err}
		}

		// Write edit buffer to connection
		_, err = connection.Write(data)

		if err != nil {
			return ConnectionError{err}
		}

		/* Response */

		header, err := getHeader(connection)

		if err != nil {
			return ConnectionError{err}
		}

		/* Validate */

		return validate(EDIT, request, header, connection)
	}

	return nil
}

// Request data deletion from the server
func (this client) delete(request REQUEST, target string, token *Token) error {
	switch this.protocol {
	case TCP:
		connection, err := this.initTCP()

		if err != nil {
			return ConnectionError{err}
		}

		defer connection.Close()

		/* Request */

		if err = setHeader(
			connection,
			DELETE,
			request,
			0,
			token,
			target,
		); err != nil {
			return ConnectionError{err}
		}

		/* Response */

		header, err := getHeader(connection)

		if err != nil {
			return ConnectionError{err}
		}

		/* Validate */

		return validate(DELETE, request, header, connection)
	}

	return nil
}

// Check if something exists on the server, without authentication
func (this client) check(request REQUEST, target string, data []byte) error {
	switch this.protocol {
	case TCP:
		connection, err := this.initTCP()

		if err != nil {
			return ConnectionError{err}
		}

		defer connection.Close()

		/* Request */

		if err = setHeader(
			connection,
			CHECK,
			request,
			uint16(len(data)),
			nil,
			target,
		); err != nil {
			return ConnectionError{err}
		}

		// Write check buffer to connection
		_, err = connection.Write(data)

		if err != nil {
			return ConnectionError{err}
		}

		/* Response */

		header, err := getHeader(connection)

		if err != nil {
			return ConnectionError{err}
		}

		/* Validate */

		return validate(CHECK, request, header, connection)
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
			return ConnectionError{err}
		}

		// Create new error
		smhbErr := errors.New(string(buffer))

		if method == QUERY {
			return NotFoundError{strconv.Itoa(int(request)), smhbErr}
		}

		return ConnectionError{smhbErr}
	}

	// Compare method
	if response.method != method {
		return ConnectionError{
			fmt.Errorf("invalid method: %d", response.method),
		}
	}

	// Compare request
	if response.request != request {
		return ConnectionError{
			fmt.Errorf("invalid response: %d", response.request),
		}
	}

	return nil
}
