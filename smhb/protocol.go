package smhb

type Header struct {
	request [2]byte
	length  [2]byte
}

type PROTOCOL int

const (
	TCP = iota
)

type REQUEST uint16

const (
	USER = REQUEST(0)
)
