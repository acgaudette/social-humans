package front

type Views map[string]interface{}

type ActiveView struct {
	Handle string
}

type StatusView struct {
	Status string
}

type UserView struct {
	Handle       string
	Name         string
	Status       string
	IsActiveUser bool
}

type PoolView struct {
	Handles []string
	Status  string
}

type FeedView struct {
	Posts  []*PostView
	Status string
}

type PostView struct {
	Title        string
	Content      string
	Author       string
	ID           string
	IsActiveUser bool
}
