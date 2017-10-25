package smhb

import (
	"encoding/binary"
	"net"
)

type PROTOCOL int

const (
	TCP = iota
)

type REQUEST uint16

const (
	USER = REQUEST(0)
)

type header struct {
	request REQUEST
	length  uint16
}

func getHeader(connection net.Conn) (header, error) {
	this := header{}

	err := binary.Read(connection, binary.LittleEndian, &this.request)

	if err != nil {
		return this, err
	}

	err = binary.Read(connection, binary.LittleEndian, &this.length)

	return this, nil
}

func setHeader(connection net.Conn, request REQUEST, length uint16) error {
	err := binary.Write(connection, binary.LittleEndian, request)

	if err != nil {
		return err
	}

	err = binary.Write(connection, binary.LittleEndian, length)

	return err
}
