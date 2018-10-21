package framework

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

const (
	pattern      = "^templates\\/.+\\.gohtml$"
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
		r.initTemplates()
	}
	err := r.templates.ExecuteTemplate(w, templateName, p)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func (r *renderer) initTemplates() {
	r.templates = template.New(templateName)
	r.templates.Funcs(map[string]interface{}{
		"RenderTemplate": r.renderTemplate,
	})
	log.Println("Loading templates...")
	wd, _ := os.Getwd()
	log.Printf("Current Working Directory: %s\n", wd)
	for _, t := range r.fp.List() {
		log.Printf("Found file: %s", t)
		if matcher.MatchString(t) {
			log.Print(" -- loading template -- ")
			r.templates.Parse(r.fp.String(t))
			log.Print("LOADED")
		}
		log.Println()
	}
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

	r.initTemplates()
	return r
}
