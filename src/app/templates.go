package app

import (
	"../views"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type templateLookup map[string]*template.Template

// Global template lookup map
var templates templateLookup

// Parse templates on initialization
func init() {
	// Sanity check
	if templates == nil {
		templates = make(templateLookup)
	}

	// Get template files
	files, err := filepath.Glob(ROOT + "/*.html")

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		base := filepath.Base(file)

		// Append template at file to base template
		targets := append([]string{ROOT + "/layout.tmpl"}, file)

		// Parse
		templates[base] = template.Must(template.ParseFiles(targets...))
	}

	log.Printf("Processed %v templates", len(files))
}

// Serve a template over HTTP
func ServeTemplate(
	out http.ResponseWriter, path string, data *views.Container,
) *Error {
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

	// Guaranteed to serve a response
	out.Header().Set("Content-Type", "text/html; charset=utf-8")
	out.Write(buffer.Bytes())

	return nil
}
