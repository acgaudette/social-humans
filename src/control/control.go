package control

import (
	"../data"
	"../front"
)

func GetUserView(user *data.User) *front.UserView {
	handle := user.Handle

	if handle == "" {
		handle = "Username Invalid"
	}

	name := user.Name

	if name == "" {
		name = "Name Invalid"
	}

	return &front.UserView{
		Handle: handle,
		Name:   name,
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
