package main

import (
	"html/template"
	"net/http"
)

type statusMessage struct {
	Status string
}

type poolUsers struct {
	Handles []string
	Status  string
}

func getPoolUsers(handle string, status string) (*poolUsers, error) {
	data, err := loadPool(handle)

	if err != nil {
		empty := &poolUsers{
			Handles: []string{},
			Status:  "Access failure",
		}

		return empty, err
	}

	if len(data.users) <= 1 {
		if status == "" {
			status = "Your pool is empty!"
		}

		empty := &poolUsers{
			Handles: []string{},
			Status:  status,
		}

		return empty, nil
	}

	result := &poolUsers{
		Handles: []string{},
		Status:  status,
	}

	for _, value := range data.users {
		if value == handle {
			continue
		}

		result.Handles = append(result.Handles, value)
	}

	return result, nil
}

func serveTemplate(
	writer http.ResponseWriter, path string, data interface{},
) error {
	t, err := template.ParseFiles(ROOT + path)

	if err != nil {
		return err
	}

	err = t.Execute(writer, data)

	return err
}
