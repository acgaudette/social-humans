package front

type LoginView struct {
	Handle     string
	Status     string
	IsLoggedIn bool
}

type UserView struct {
	Handle       string
	Name         string
	Status       string
	IsActiveUser bool
}

type PoolView struct {
	Handle  string
	Handles []string
	Status  string
}

type PostView struct {
	Handle string
	Status string
}
