package smhb

import (
	"fmt"
)

type NotFoundError struct {
	identifier string
	err        error
}

func (this NotFoundError) Error() string {
	return fmt.Sprintf("%s not found: %s", this.identifier, this.err)
}

type ConnectionError struct {
	err error
}

func (this ConnectionError) Error() string {
	return fmt.Sprintf("error communicating with server: %s", this.err)
}
