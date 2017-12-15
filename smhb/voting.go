package smhb

import (
	"net"
	"time"
)

func connect(destination string) (net.Conn, error) {
	connection, err := net.DialTimeout("tcp", destination, time.Second*20)

	if err != nil {
		return nil, err
	}

	return connection, nil
}

func proposeTransaction(
	request REQUEST,
	target string,
	data []byte,
	timestamp string,
	destination string,
) error {
	connection, err := connect(destination)

	if err != nil {
		return ConnectionError{err}
	}

	defer connection.Close()

	/* Request */

	// No token checking for replication processes (RIP)
	token := Token{}

	if err = setHeader(
		connection,
		PROPOSE,
		request,
		uint16(len(data)),
		&token,
		target,
	); err != nil {
		return ConnectionError{err}
	}

	// Write store buffer to connection
	_, err = connection.Write(data)

	if err != nil {
		return ConnectionError{err}
	}

	return nil
}

func ackTransaction(timestamp string) error {
	// TODO: implement
	return nil
}

func commitTransaction(timestamp string) error {
	// TODO: implement
	// delete the vote from the map!
	return nil
}
