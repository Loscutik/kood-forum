package templates

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"forum/app/config"
)

const (
	TEMPLATES_PATH = "./webui/templates/"
	STATIC_PATH    = "./webui/static/"
)

/*
returnes all parsed templates
*/
func NewTemplateCache(templateDir string) (map[string]*template.Template, error) {
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
func ExecuteTemplate(app *config.Application, w http.ResponseWriter, r *http.Request, name string, outputData any) error {
	tm, ok := app.TemlateCashe[name]
	if !ok {
		return fmt.Errorf("the template '%s' is not found", name)
	}

	err := tm.Execute(w, outputData)
	if err != nil {
		return fmt.Errorf("the template '%s' executing is failed: %v", name, err)
	}
	return nil
}

func ExecuteError(app *config.Application, w http.ResponseWriter, r *http.Request, statusCode int) {
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
		app.ErrLog.Printf("can't parse %s template: %v", pageName, err)
		http.Error(w, fmt.Sprintf("ERROR: %s. ", http.StatusText(statusCode)), statusCode)
	}
	err = tm.Execute(w, nil)
	if err != nil {
		app.ErrLog.Printf("can't execute  %s template: %v", pageName, err)
		http.Error(w, fmt.Sprintf("ERROR: %s. ", http.StatusText(statusCode)), statusCode)
	}
}
