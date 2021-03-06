package smhb

import (
	"fmt"
)

/* Data not found on server */

type NotFoundError struct {
	identifier string
	err        error
}

func (this NotFoundError) Error() string {
	return fmt.Sprintf("%s not found: %s", this.identifier, this.err)
}

/* Error during client/server communication */

type ConnectionError struct {
	err error
}

func (this ConnectionError) Error() string {
	return fmt.Sprintf("error communicating with server: %s", this.err)
}

/* Invalid credentials */

type AuthError struct {
	err error
}

func (this AuthError) Error() string {
	return fmt.Sprintf("error authenticating with server: %s", this.err)
}
