package front

type StatusMessage struct {
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
