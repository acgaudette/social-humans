package front

import (
	"../app"
	"bytes"
	"fmt"
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
		targets := append([]string{ROOT + "/layout.tmpl"}, file)
		templates[base] = template.Must(template.ParseFiles(targets...))
	}

	log.Printf("Processed %v templates", len(files))
}

// Guaranteed to serve a response
func ServeTemplate(
	out http.ResponseWriter, path string, data *Views,
) *app.Error {
	// Load template from cache
	target, ok := templates[path+".html"]

	if !ok {
		return ServerError(
			fmt.Errorf("template \"%s\" does not exist", path),
		)
	}

	// Write template output to buffer to prevent a dirty response
	var buffer bytes.Buffer
	err := target.ExecuteTemplate(&buffer, "layout", data)

	if err != nil {
		return ServerError(err)
	}

	out.Header().Set("Content-Type", "text/html; charset=utf-8")
	out.Write(buffer.Bytes())

	return nil
}
