package control

import (
	"../data"
	"../front"
	"log"
	"net/http"
)

func GetUserView(
	user *data.User, status string, request *http.Request,
) *front.UserView {
	handle := user.Handle

	if handle == "" {
		handle = "Username Invalid"
	}

	name := user.Name

	if name == "" {
		name = "Name Invalid"
	}

	account, err := data.GetUserFromSession(request)

	if err != nil {
		log.Printf("%s", err)
	}

	active := account.Handle == user.Handle

	return &front.UserView{
		Handle:       handle,
		Name:         name,
		Status:       status,
		IsActiveUser: active,
	}
}

func GetPoolView(handle string, status string) (*front.PoolView, error) {
	pool, err := data.LoadPool(handle)

	if err != nil {
		empty := &front.PoolView{
			Handles: []string{},
			Status:  "Error: access failure",
		}

		return empty, err
	}

	if len(pool.Users) <= 1 {
		if status == "" {
			status = "Your pool is empty!"
		}

		empty := &front.PoolView{
			Handles: []string{},
			Status:  status,
		}

		return empty, nil
	}

	result := &front.PoolView{
		Handle:  handle,
		Handles: []string{},
		Status:  status,
	}

	for _, value := range pool.Users {
		if value == handle {
			continue
		}

		result.Handles = append(result.Handles, value)
	}

	return result, nil
}
