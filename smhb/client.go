package smhb

import (
	"errors"
	"io"
	"net"
	"strconv"
)

type Client interface {
	ServerAddress() string
	ServerPort() int
	Protocol() PROTOCOL

	GetUser(string) (User, error)
	GetPool(string) (Pool, error)
	GetPost(string) (Post, error)
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

		err = setHeader(connection, request, 0, target)

		if err != nil {
			return nil, err
		}

		// Response

		header, err := getHeader(connection)

		if err != nil {
			return nil, err
		}

		// Validate

		if header.request != request {
			return nil, errors.New("invalid response")
		}

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
