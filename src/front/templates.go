package front

import (
	"bytes"
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

	log.Printf("Processed %v templates", len(files))
}

func ServeTemplate(
	writer http.ResponseWriter, path string, data interface{},
) error {
	target, ok := templates[path+".html"]

	if !ok {
		err := errors.New("template does not exist")
		Error501(writer)
		return err
	}

	var buffer bytes.Buffer
	err := target.ExecuteTemplate(&buffer, "layout", data)

	if err != nil {
		Error501(writer)
		return err
	}

	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.Write(buffer.Bytes())

	return nil
}
