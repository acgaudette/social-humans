package smhb

import (
	"encoding/binary"
	"net"
	"strconv"
)

type Client interface {
	ServerAddress() string
	ServerPort() int
	Protocol() PROTOCOL
	Query(request REQUEST) error
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

func (this client) Query(request REQUEST) error {
	switch this.protocol {
	case TCP:
		bind := this.serverAddress + ":" + strconv.Itoa(this.serverPort)
		connection, err := net.Dial("tcp", bind)

		if err != nil {
			return err
		}

		defer connection.Close()

		err = binary.Write(connection, binary.LittleEndian, request)

		if err != nil {
			return err
		}

		err = binary.Write(connection, binary.LittleEndian, uint16(4)) // fix this later

		if err != nil {
			return err
		}
	}

	return nil
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
