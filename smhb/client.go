package smhb

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strconv"
)

type Client interface {
	ServerAddress() string
	ServerPort() int
	Protocol() PROTOCOL
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

func (this client) query(request REQUEST) ([]byte, error) {
	switch this.protocol {
	case TCP:
		bind := this.serverAddress + ":" + strconv.Itoa(this.serverPort)
		connection, err := net.Dial("tcp", bind)

		if err != nil {
			return nil, err
		}

		defer connection.Close()

		// Request

		err = binary.Write(connection, binary.LittleEndian, request)

		if err != nil {
			return nil, err
		}

		err = binary.Write(connection, binary.LittleEndian, uint16(4))

		if err != nil {
			return nil, err
		}

		// Response

		var response, length uint16

		err = binary.Read(connection, binary.LittleEndian, &request)

		if err != nil {
			return nil, err
		}

		err = binary.Read(connection, binary.LittleEndian, &length)

		if err != nil {
			return nil, err
		}

		// Validate

		if REQUEST(response) != request {
			return nil, errors.New("invalid response")
		}

		length -= 4
		buffer := make([]byte, length)
		_, err = io.ReadFull(connection, buffer)

		if err != nil {
			return nil, err
		}

		return buffer, nil
	}

	return nil, nil
}
