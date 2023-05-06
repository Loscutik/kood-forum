package templates

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

type TemplateCache map[string]*template.Template

/*
returnes all parsed templates
*/
func NewTemplateCache(templateDir string) (TemplateCache, error) {
	temlateCashe := TemplateCache{}
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
func ExecuteTemplate(templateCache TemplateCache, w http.ResponseWriter, r *http.Request, name string, outputData any) error {
	tm, ok := templateCache[name]
	if !ok {
		return fmt.Errorf("the template '%s' is not found", name)
	}

	err := tm.Execute(w, outputData)
	if err != nil {
		return fmt.Errorf("the template '%s' executing is failed: %v", name, err)
	}
	return nil
}

/*
executes a template for the given error (statusCode)
*/
func ExecuteError(w http.ResponseWriter, r *http.Request, statusCode int) error {
	var pageName string
	switch statusCode {
	case http.StatusNotFound:
		pageName = "error404.html"
	case http.StatusForbidden:
		pageName = "forbidden.tmpl"
	default:
		pageName = "error404.html"
	}

	tm, err := template.ParseFiles(TEMPLATES_PATH+pageName, TEMPLATES_PATH+"base.layout.tmpl") // Opens the HTML web page
	if err != nil {
		return fmt.Errorf("can't parse %s template: %v", pageName, err)
	}
	err = tm.Execute(w, nil)
	if err != nil {
		return fmt.Errorf("can't execute  %s template: %v", pageName, err)
	}
	return nil
}
