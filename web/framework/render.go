package framework

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

const (
	pattern      = `^templates(\\|\/)\w+\.gohtml$` // matches templates/*.gohtml (Linux & macOS) or templates\*.gohtml (Windows)
	templateName = "master"
)

var matcher *regexp.Regexp

func init() {
	matcher = regexp.MustCompile(pattern)
}

// Renderer type
type Renderer interface {
	Render(w http.ResponseWriter, p ViewModel)
}

// ViewModel defines the page ViewModel
type ViewModel struct {
	Template string
	Title    string
	Data     interface{}
}

type renderer struct {
	fp        FileProvider
	templates *template.Template
	config    *Config
}

// Render the specified view model
func (r *renderer) Render(w http.ResponseWriter, p ViewModel) {
	if r.config.Debug {
		err := r.initTemplates()
		if err != nil {
			log.Printf(err.Error())
		}
	}
	err := r.templates.ExecuteTemplate(w, templateName, p)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func (r *renderer) initTemplates() error {
	r.templates = template.New(templateName)
	r.templates.Funcs(map[string]interface{}{
		"RenderTemplate": r.renderTemplate,
	})
	for _, t := range r.fp.List() {
		if matcher.MatchString(t) {
			_, err := r.templates.Parse(r.fp.String(t))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *renderer) renderTemplate(name string, data interface{}) (ret template.HTML, err error) {
	buf := bytes.NewBuffer([]byte{})
	err = r.templates.ExecuteTemplate(buf, name, data)
	ret = template.HTML(buf.String())
	return
}

// NewRenderer creates a new rendere
func NewRenderer(fp FileProvider, config *Config) Renderer {
	r := &renderer{
		config: config,
		fp:     fp,
	}

	err := r.initTemplates()
	if err != nil {
		log.Fatalln(err.Error())
	}
	return r
}
