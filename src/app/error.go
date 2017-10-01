package app

type ERROR_CODE int

const (
	SERVER = iota
	NOT_FOUND
	REDIRECT
)

type Error struct {
	Native error
	Code   ERROR_CODE
}
