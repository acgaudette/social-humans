package main

import(
	"net/http"
	"html/template"
)

type statusMessage struct {
	Status string
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
