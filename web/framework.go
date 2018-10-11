package framework

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/bashmohandes/go-askme/shared"
	"github.com/julienschmidt/httprouter"
)

// Controller represents Controller in MVC  model
type Controller struct {
	templates *template.Template
	actions   []*Action
	fp        shared.FileProvider
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

// Init initializes the controller
func (c *Controller) Init(fp shared.FileProvider) {
	c.templates = template.New("master")
	c.actions = make([]*Action, 0)
	c.fp = fp
	c.templates.Funcs(map[string]interface{}{
		"RenderTemplate": c.renderTemplate,
	})
	c.loadTemplates()
}

// AddAction adds action to controller
func (c *Controller) AddAction(method string, path string, f httprouter.Handle) {
	c.actions = append(c.actions, &Action{Method: method, Path: path, Func: f})
}

// Render the specified view model
func (c *Controller) Render(w http.ResponseWriter, p ViewModel) {
	err := c.templates.ExecuteTemplate(w, "master", p)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

// Actions return the list of configured actions
func (c *Controller) Actions() []*Action {
	return c.actions
}

func (c *Controller) loadTemplates() {
	for _, t := range c.fp.List() {
		if strings.HasSuffix(t, ".gohtml") {
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
