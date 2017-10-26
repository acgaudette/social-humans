package smhb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
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
	POOL           = REQUEST(1)
	POOL_ADD       = REQUEST(2)
	POOL_BLOCK     = REQUEST(3)
	POST_ADDRESSES = REQUEST(4)
	POST           = REQUEST(5)
)

type header struct {
	method  METHOD
	request REQUEST
	length  uint16
	target  string
}

func getHeader(connection net.Conn) (header, error) {
	this := header{}

	err := binary.Read(connection, binary.LittleEndian, &this.method)

	if err != nil {
		return this, err
	}

	err = binary.Read(connection, binary.LittleEndian, &this.request)

	if err != nil {
		return this, err
	}

	err = binary.Read(connection, binary.LittleEndian, &this.length)

	if err != nil {
		return this, err
	}

	var buffer [TARGET_LENGTH]byte
	_, err = connection.Read(buffer[:])
	end := bytes.IndexByte(buffer[:], byte('\000'))

	if end < 0 {
		return this, err
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
	err := binary.Write(connection, binary.LittleEndian, method)

	if err != nil {
		return err
	}

	err = binary.Write(connection, binary.LittleEndian, request)

	if err != nil {
		return err
	}

	err = binary.Write(connection, binary.LittleEndian, length)

	if err != nil {
		return err
	}

	var buffer [TARGET_LENGTH]byte
	copied := copy(buffer[:], target)

	if copied < len(target) {
		return errors.New("target string overflow")
	}

	_, err = connection.Write(buffer[:])

	return err
}
