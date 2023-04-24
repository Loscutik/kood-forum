package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	TEMPLATES_PATH = "./webui/templates/"
	STATIC_PATH    = "./webui/static/"
)

/*
returnes all parsed templates
*/
func newTemplateCache(templateDir string) (map[string]*template.Template, error) {
	temlateCashe := map[string]*template.Template{}
	// get all templates of pages

	pages, err := filepath.Glob(filepath.Join(templateDir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}
	// create templates for all pages
	for _, page := range pages {
		tm, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// add a layout template to the each page
		tm, err = tm.ParseGlob(filepath.Join(templateDir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// add partial templates to the each page
		tm, err = tm.ParseGlob(filepath.Join(templateDir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		
		temlateCashe[strings.TrimSuffix(tm.Name(), ".page.tmpl")] = tm
	}
	return temlateCashe, nil
}

/*
executes a template with the given name using the given data
*/
func (app *application) executeTemplate(w http.ResponseWriter, r *http.Request, name string, outputData any) {
	tm, ok := app.temlateCashe[name]
	if !ok {
		app.ServerError(w, r, fmt.Sprintf("the template '%s' is not found", name), nil)
		return
	}

	err := tm.Execute(w, outputData)
	if err != nil {
		app.ServerError(w, r, "the template executing is failed", err)
		return
	}
}
