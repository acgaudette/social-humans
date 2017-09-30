package front

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type templateLookup map[string]*template.Template

var templates templateLookup

func init() {
	if templates == nil {
		templates = make(templateLookup)
	}

	files, err := filepath.Glob(ROOT + "/*.html")

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		base := filepath.Base(file)

		if base == "layout.html" {
			continue
		}

		targets := append([]string{ROOT + "/layout.html"}, file)
		templates[base] = template.Must(template.ParseFiles(targets...))
	}
}

func ServeTemplate(
	writer http.ResponseWriter, path string, data interface{},
) error {
	target, ok := templates[path+".html"]

	if !ok {
		return errors.New("template does not exist")
	}

	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	return target.ExecuteTemplate(writer, "layout", data)
}
