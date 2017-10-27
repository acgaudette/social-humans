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
