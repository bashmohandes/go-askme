package framework

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"regexp"

	"github.com/bashmohandes/go-askme/shared"
	"github.com/julienschmidt/httprouter"
)

// Controller represents Controller in MVC  model
type Controller struct {
	templates *template.Template
	actions   []*Action
	fp        shared.FileProvider
	config    *Config
}

// Action captures the http actions of the controller
type Action struct {
	Method string
	Path   string
	Func   httprouter.Handle
}

// ViewModel defines the page ViewModel
type ViewModel struct {
	Template string
	Title    string
	Data     interface{}
}

// Config configuration
type Config struct {
	Debug bool
	// Assets relative path to askme package
	Assets string
	Port   int
}

const (
	pattern      = "^templates\\/.+\\.gohtml$"
	templateName = "master"
)

var matcher *regexp.Regexp

func init() {
	matcher = regexp.MustCompile(pattern)
}

// Init initializes the controller
func (c *Controller) Init(fp shared.FileProvider, config *Config) {
	c.config = config
	c.actions = make([]*Action, 0)
	c.fp = fp
	c.initTemplates()
}

// AddAction adds action to controller
func (c *Controller) AddAction(method string, path string, f httprouter.Handle) {
	c.actions = append(c.actions, &Action{method, path, f})
}

// Render the specified view model
func (c *Controller) Render(w http.ResponseWriter, p ViewModel) {
	if c.config.Debug {
		c.initTemplates()
	}
	err := c.templates.ExecuteTemplate(w, templateName, p)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

// Actions return the list of configured actions
func (c *Controller) Actions() []*Action {
	return c.actions
}

func (c *Controller) initTemplates() {
	c.templates = template.New(templateName)
	c.templates.Funcs(map[string]interface{}{
		"RenderTemplate": c.renderTemplate,
	})
	for _, t := range c.fp.List() {
		if matcher.MatchString(t) {
			c.templates.Parse(c.fp.String(t))
		}
	}
}

func (c *Controller) renderTemplate(name string, data interface{}) (ret template.HTML, err error) {
	buf := bytes.NewBuffer([]byte{})
	err = c.templates.ExecuteTemplate(buf, name, data)
	ret = template.HTML(buf.String())
	return
}
