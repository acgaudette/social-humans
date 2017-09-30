package control

import (
	"../data"
	"../front"
)

func GetUserView(user *data.User) *front.UserView {
	return &front.UserView{
		Handle: user.Handle,
	}
}

func GetPoolView(handle string, status string) (*front.PoolView, error) {
	pool, err := data.LoadPool(handle)

	if err != nil {
		empty := &front.PoolView{
			Handles: []string{},
			Status:  "Access failure",
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
