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
	IsActiveUser bool
}

type PoolView struct {
	Handles []string
}

type FeedView struct {
	Posts []*PostView
}

type PostView struct {
	Title        string
	Content      string
	Author       string
	ID           string
	Timestamp    string
	IsActiveUser bool
}
