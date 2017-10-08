package control

import "fmt"

type EmptyFeedError struct {
	handle string
}

func (this EmptyFeedError) Error() string {
	return fmt.Sprintf("feed for user \"%s\" is empty", this.handle)
}

type EmptyPoolError struct {
	handle string
}

func (this EmptyPoolError) Error() string {
	return fmt.Sprintf("pool for user \"%s\" is empty", this.handle)
}

type AccessError struct {
	handle string
}

func (this AccessError) Error() string {
	return fmt.Sprintf("feed for user \"%s\" access failure", this.handle)
}
