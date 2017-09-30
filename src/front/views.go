package front

type StatusMessage struct {
	Status string
}

type UserView struct {
	Handle string
}

type PoolView struct {
	Handles []string
	Status  string
}
