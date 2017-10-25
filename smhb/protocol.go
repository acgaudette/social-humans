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

type REQUEST int16

const (
	ERROR          = REQUEST(-1)
	USER           = REQUEST(0)
	POOL           = REQUEST(1)
	POST           = REQUEST(2)
	POST_ADDRESSES = REQUEST(3)
)

type header struct {
	request REQUEST
	length  uint16
	target  string
}

func getHeader(connection net.Conn) (header, error) {
	this := header{}

	err := binary.Read(connection, binary.LittleEndian, &this.request)

	if err != nil {
		return this, err
	}

	err = binary.Read(connection, binary.LittleEndian, &this.length)

	if err != nil {
		return this, err
	}

	var buffer [28]byte
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
	request REQUEST,
	length uint16,
	target string,
) error {
	err := binary.Write(connection, binary.LittleEndian, request)

	if err != nil {
		return err
	}

	err = binary.Write(connection, binary.LittleEndian, length)

	if err != nil {
		return err
	}

	var buffer [28]byte
	copied := copy(buffer[:], target)

	if copied < len(target) {
		return errors.New("target string overflow")
	}

	_, err = connection.Write(buffer[:])

	return err
}
