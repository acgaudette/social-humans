package views

// Type alias for readability
type View interface{}

/* Views */

type Active struct {
	Handle string
}

type Status struct {
	Status string
}

type User struct {
	Handle       string
	Name         string
	IsActiveUser bool
}

type Pool struct {
	Handles []string
}

type Feed struct {
	Posts []*Post
}

type Post struct {
	Title        string
	Content      string
	Author       string
	ID           string
	WasEdited    bool
	Timestamp    string
	IsActiveUser bool
}

type Base struct {
	Commit string
	Link   string
}
