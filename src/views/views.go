package views

type View interface{}

type Container map[string]View

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
	Timestamp    string
	IsActiveUser bool
}
