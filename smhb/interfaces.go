package smhb

/*
	Models are implemented with interfaces so that the structures remain bound
	to real data
*/

// Public-facing user interface
type User interface {
	Handle() string
	Name() string
	Equals(User) bool
}

// Public-facing pool interface
type Pool interface {
	Handle() string
	Users() userPool
}

// Public-facing post interface
type Post interface {
	Title() string
	Content() string
	Author() string
	Timestamp() string
	WasEdited() bool
	WasAuthoredBy(string) bool
	GetAddress() string
}

// Public-facing feed interface
type Feed interface {
	Addresses() []string
}
