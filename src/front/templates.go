package front

import (
	"../app"
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

// Guaranteed to serve a response
func ServeTemplate(
	out http.ResponseWriter, path string, data interface{},
) *app.Error {
	// Load template from cache
	target, ok := templates[path+".html"]

	if !ok {
		return &app.Error{
			Native: errors.New("template does not exist"),
			Code:   app.SERVER,
		}
	}

	// Write template output to buffer to prevent a dirty response
	var buffer bytes.Buffer
	err := target.ExecuteTemplate(&buffer, "layout", data)

	if err != nil {
		return &app.Error{
			Native: err,
			Code:   app.SERVER,
		}
	}

	out.Header().Set("Content-Type", "text/html; charset=utf-8")
	out.Write(buffer.Bytes())

	return nil
}
