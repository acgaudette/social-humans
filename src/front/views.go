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
