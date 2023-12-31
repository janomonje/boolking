package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/janomonje/booking/pkg/config"
	"github.com/janomonje/booking/pkg/models"
)

/*
// Renders template using html/template
func RenderTemplateTest(w http.ResponseWriter, tmpl string) {
	parsedTemplate, _ := template.ParseFiles("./cmd/templates/"+tmpl, "./cmd/templates/base.layout.tmpl")
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println("error passing template:", err)
	}
}

// Stores rendered templates
var tc = make(map[string]*template.Template)

/*
*****

	This function will take two variables, a response writer and the template we want to render.
	 *****

func RenderTemplate(w http.ResponseWriter, t string) {
	var tmpl *template.Template
	var err error

	// check to see if we already have the template in our cache
	_, inMap := tc[t]

	if !inMap {
		// need to create the template
		log.Println("Creating template and adding to cache")
		err = createTemplateCache(t)
		if err != nil {
			log.Println(err)
		}
	} else {
		// we have the template in the cache
		log.Println("using cached template")
	}

	tmpl = tc[t]

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}

}

// Populating "tmpl = tc[t]"
func createTemplateCache(t string) error {
	templates := []string{
		fmt.Sprintf("./cmd/templates/%s", t), "./cmd/templates/base.layout.tmpl",
	}

	// parsing the template
	tmpl, err := template.ParseFiles(templates...) // takes the entries from tamplates and places them as individual strings

	if err != nil {
		return err
	}
	// add template to cache(map)
	tc[t] = tmpl

	return nil
}
*/

// *** More complex render cache ***
var app *config.AppConfig

// sets the config for the config package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// Renders template using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {

	var templateCache map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	// get requested template from cache
	template, ok := templateCache[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td)
	_ = template.Execute(buf, td)

	// render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

}

func CreateTemplateCache() (map[string]*template.Template, error) {
	//myCache := make(map[string]template.Template)
	myCache := map[string]*template.Template{}

	// get all of the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// range thorugh all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page) // ts= template set

		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
