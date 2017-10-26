package smhb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
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
)

type REQUEST int16

const (
	ERROR          = REQUEST(-1)
	USER           = REQUEST(0)
	USER_NAME      = REQUEST(1)
	USER_PASSWORD  = REQUEST(2)
	POOL           = REQUEST(3)
	POOL_ADD       = REQUEST(4)
	POOL_BLOCK     = REQUEST(5)
	POST_ADDRESSES = REQUEST(6)
	POST           = REQUEST(7)
)

type header struct {
	method  METHOD
	request REQUEST
	length  uint16
	target  string
}

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

	// Read target string

	var buffer [TARGET_LENGTH]byte
	_, err = connection.Read(buffer[:])

	if err != nil {
		return this, err
	}

	end := bytes.IndexByte(buffer[:], byte('\000'))

	if end < 0 {
		return this, errors.New("target string not terminated")
	}

	this.target = string(buffer[:end])

	return this, nil
}

func setHeader(
	connection net.Conn,
	method METHOD,
	request REQUEST,
	length uint16,
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

	// Write target string

	var buffer [TARGET_LENGTH]byte
	copied := copy(buffer[:], target)

	if copied < len(target) {
		return errors.New("target string overflow")
	}

	_, err = connection.Write(buffer[:])

	return err
}
