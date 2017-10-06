package app

type ERROR_CODE int

const (
	SERVER = iota
	NOT_FOUND
	FORBIDDEN
	REDIRECT
)

// App error passed to handler
type Error struct {
	Native error
	Code   ERROR_CODE
}
