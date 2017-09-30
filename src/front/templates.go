package front

import (
	"../data"
	"html/template"
	"net/http"
)

type StatusMessage struct {
	Status string
}

type PoolUsers struct {
	Handles []string
	Status  string
}

func GetPoolUsers(handle string, status string) (*PoolUsers, error) {
	data, err := data.LoadPool(handle)

	if err != nil {
		empty := &PoolUsers{
			Handles: []string{},
			Status:  "Access failure",
		}

		return empty, err
	}

	if len(data.Users) <= 1 {
		if status == "" {
			status = "Your pool is empty!"
		}

		empty := &PoolUsers{
			Handles: []string{},
			Status:  status,
		}

		return empty, nil
	}

	result := &PoolUsers{
		Handles: []string{},
		Status:  status,
	}

	for _, value := range data.Users {
		if value == handle {
			continue
		}

		result.Handles = append(result.Handles, value)
	}

	return result, nil
}

func ServeTemplate(
	writer http.ResponseWriter, path string, data interface{},
) error {
	t, err := template.ParseFiles(ROOT + path)

	if err != nil {
		return err
	}

	err = t.Execute(writer, data)

	return err
}
