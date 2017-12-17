package smhb

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

// NTP epoch
func epoch() time.Time {
	return time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
}

type fixed64 uint64

func (this fixed64) toDuration() time.Duration {
	// Convert upper 32 bits to time.Duration
	nanoseconds := time.Duration(this>>32) * time.Second
	// Multiply fraction by smallest unit (time.Duration) and divide
	remainder := time.Duration(this&0xffffffff) * time.Second >> 32
	return nanoseconds + remainder
}

func (this fixed64) toTime() time.Time {
	return epoch().Add(this.toDuration())
}

func toFixed64(this time.Time) fixed64 {
	nanoseconds := this.Sub(epoch())
	seconds := nanoseconds / time.Second
	remainder := nanoseconds - seconds*time.Second
	// Multiply remainder and correct for unit (time.Duration)
	fraction := (remainder << 32) / time.Second
	return fixed64(seconds<<32 | fraction)
}

// NTP packet
type ntp struct {
	LeapVersionMode uint8
	Stratum         uint8
	Poll            int8
	Precision       int8
	RootDelay       uint32
	RootDispersion  uint32
	RefID           uint32
	RefTime         fixed64
	OriginTime      fixed64
	ReceiveTime     fixed64
	TransmitTime    fixed64
}

// Calculate the current time via a remote NTP server
func getNTPTime() (*time.Time, error) {
	// Resolve NTP server hostname
	address, err := net.ResolveUDPAddr(
		"udp",
		net.JoinHostPort(NTP_SERVER, "123"),
	)

	if err != nil {
		return nil, err
	}

	// Connect
	connection, err := net.DialUDP("udp", nil, address)

	if err != nil {
		return nil, err
	}

	defer connection.Close()

	// Set timeout
	connection.SetDeadline(time.Now().Add(time.Second * NTP_TIMEOUT))

	// Build packets
	query := new(ntp)
	response := new(ntp)

	query.LeapVersionMode = 227 // 0b11100011
	sent := time.Now()
	query.TransmitTime = toFixed64(sent)

	// Query
	err = binary.Write(connection, binary.BigEndian, query)

	if err != nil {
		return nil, fmt.Errorf("error writing to connection: %s", err)
	}

	// Response
	err = binary.Read(connection, binary.BigEndian, response)

	if err != nil {
		return nil, fmt.Errorf("error reading from connection: %s", err)
	}

	delta := time.Since(sent)
	receive := sent.Add(delta)

	// Calculate offset using four timestamps
	l := response.ReceiveTime.toTime().Sub(sent)
	r := response.TransmitTime.toTime().Sub(receive)
	offset := (l + r) / 2

	result := time.Now().Add(offset)
	return &result, nil
}
