package front

import (
	"html/template"
	"net/http"
)

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
