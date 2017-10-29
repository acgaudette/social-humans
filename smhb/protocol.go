package smhb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"fmt"
	"time"
)

type PROTOCOL int

const (
	TCP = iota
)

type METHOD uint16

const (
	QUERY  = METHOD(0)
	STORE  = METHOD(1)
	EDIT   = METHOD(2)
	DELETE = METHOD(3)
	CHECK  = METHOD(4)
)

type REQUEST int16

const (
	ERR_AUTH       = REQUEST(-3)
	ERR_NOT_FOUND  = REQUEST(-2)
	ERR            = REQUEST(-1)
	VALIDATE       = REQUEST(0)
	TOKEN          = REQUEST(1)
	USER           = REQUEST(2)
	USER_NAME      = REQUEST(3)
	USER_PASSWORD  = REQUEST(4)
	POOL           = REQUEST(5)
	POOL_ADD       = REQUEST(6)
	POOL_BLOCK     = REQUEST(7)
	POST_ADDRESSES = REQUEST(8)
	POST           = REQUEST(9)
	FEED           = REQUEST(10)
)

// Protocol header
type header struct {
	method  METHOD  // 2 bytes
	request REQUEST // 2 bytes
	length  uint16  // 2 bytes
	token   Token   // TOKEN_SIZE bytes
	target  string  // TARGET_LENGTH + 1 bytes
}

// Read header from connection
func getHeader(connection net.Conn) (header, error) {
	this := header{}

	// Set timeout
	connection.SetDeadline(time.Now().Add(IO_TIMEOUT * time.Second))

	// Read method

	err := binary.Read(connection, binary.LittleEndian, &this.method)

	if err != nil {
		return this, err
	}

	// Read request

	err = binary.Read(connection, binary.LittleEndian, &this.request)

	if err != nil {
		return this, err
	}

	// Read data length

	err = binary.Read(connection, binary.LittleEndian, &this.length)

	if err != nil {
		return this, err
	}

	// Read token

	var tokenBuffer [TOKEN_SIZE]byte
	_, err = connection.Read(tokenBuffer[:])

	if err != nil {
		return this, err
	}

	this.token = NewToken(string(tokenBuffer[:TOKEN_SIZE]))

	// Read target string

	var targetBuffer [TARGET_LENGTH + 1]byte
	_, err = connection.Read(targetBuffer[:])

	if err != nil {
		return this, err
	}

	end := bytes.IndexByte(targetBuffer[:], byte('\000'))

	if end < 0 {
		return this, errors.New(
			"target string not terminated: corrupted or overflowed string",
		)
	}

	this.target = string(targetBuffer[:end])

	return this, nil
}

// Write header to connection
func setHeader(
	connection net.Conn,
	method METHOD,
	request REQUEST,
	length uint16,
	token *Token,
	target string,
) error {
	// Set timeout
	connection.SetDeadline(time.Now().Add(IO_TIMEOUT * time.Second))

	// Write method

	err := binary.Write(connection, binary.LittleEndian, method)

	if err != nil {
		return err
	}

	// Write request

	err = binary.Write(connection, binary.LittleEndian, request)

	if err != nil {
		return err
	}

	// Write data length

	err = binary.Write(connection, binary.LittleEndian, length)

	if err != nil {
		return err
	}

	// Write token

	var tokenBuffer [TOKEN_SIZE]byte

	if token != nil {
		copied := copy(tokenBuffer[:], token.value) // Chop null-terminator

		if copied < len(token.value)-1 {
			return errors.New("token overflow")
		}
	}

	_, err = connection.Write(tokenBuffer[:])

	// Write target string

	var targetBuffer [TARGET_LENGTH + 1]byte
	copied := copy(targetBuffer[:], target)

	if copied < len(target) {
		return errors.New("target string overflow")
	}

	_, err = connection.Write(targetBuffer[:])

	return err
}

/* Debug functions */

func (this METHOD) ToString() string {
	switch this {
	case QUERY:
		return "QUERY"
	case STORE:
		return "STORE"
	case EDIT:
		return "EDIT"
	case DELETE:
		return "DELETE"
	default:
		return fmt.Sprintf("?%v", this)
	}
}

func (this REQUEST) ToString() string {
	switch this {
	case ERR_AUTH:
		return "ERR_AUTH"
	case ERR_NOT_FOUND:
		return "ERR_NOT_FOUND"
	case ERR:
		return "ERR"
	case VALIDATE:
		return "VALIDATE"
	case TOKEN:
		return "TOKEN"
	case USER:
		return "USER"
	case USER_NAME:
		return "USER_NAME"
	case USER_PASSWORD:
		return "USER_PASSWORD"
	case POOL:
		return "POOL"
	case POOL_ADD:
		return "POOL_ADD"
	case POOL_BLOCK:
		return "POOL_BLOCK"
	case POST_ADDRESSES:
		return "POST_ADDRESSES"
	case POST:
		return "POST"
	case FEED:
		return "FEED"
	default:
		return fmt.Sprintf("?%v", this)
	}
}
