package smhb

/*
	Models are implemented with interfaces so that the structures remain bound
	to real data
*/

// Public-facing user interface
type User interface {
	Handle() string
	Name() string
	SetName(string) error
	Validate(string) error
	UpdatePassword(string) error
	Equals(User) bool
}

// Store user pool as a set of handles
type userPool map[string]string

// Public-facing pool interface
type Pool interface {
	Handle() string
	Users() userPool
	Add(string) error
	Block(string) error
}

// Public-facing post interface
type Post interface {
	Title() string
	Content() string
	Author() string
	Timestamp() string
	WasEdited() bool
	Update(string, string) error
	WasAuthoredBy(string) bool
}
