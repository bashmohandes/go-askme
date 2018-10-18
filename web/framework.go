package framework

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/bashmohandes/go-askme/shared"
	"github.com/julienschmidt/httprouter"
)

// Router interface
type Router interface {
	Actions() []*Action
	Get(path string, f httprouter.Handle)
	Post(path string, f httprouter.Handle)
	Delete(path string, f httprouter.Handle)
	Put(path string, f httprouter.Handle)
}

// Renderer type
type Renderer interface {
	Render(w http.ResponseWriter, p ViewModel)
}

// App represents the AskMe application server
type App struct {
	config       *Config
	fileProvider shared.FileProvider
	rtr          Router
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
	Debug        bool
	Port         int
	PublicFolder string
}

// router represents router in MVC  model
type router struct {
	actions []*Action
}

type renderer struct {
	fp        shared.FileProvider
	templates *template.Template
	config    *Config
}

const (
	pattern      = "^templates\\/.+\\.gohtml$"
	templateName = "master"
)

var matcher *regexp.Regexp

func init() {
	matcher = regexp.MustCompile(pattern)
}

// NewRouter initializes the controller
func NewRouter() Router {
	return &router{
		actions: make([]*Action, 0),
	}
}

// Actions return the list of configured actions
func (r *router) Actions() []*Action {
	return r.actions
}

func (r *router) Get(path string, f httprouter.Handle) {
	r.action("GET", path, f)
}

func (r *router) Post(path string, f httprouter.Handle) {
	r.action("POST", path, f)
}

func (r *router) Delete(path string, f httprouter.Handle) {
	r.action("DELETE", path, f)
}

func (r *router) Put(path string, f httprouter.Handle) {
	r.action("PUT", path, f)
}

// AddAction adds action to controller
func (r *router) action(method string, path string, f httprouter.Handle) {
	log.Printf("caller %T, method %s, path %s\n", r, method, path)
	r.actions = append(r.actions, &Action{method, path, f})
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
func NewRenderer(fp shared.FileProvider, config *Config) Renderer {
	r := &renderer{
		config: config,
		fp:     fp,
	}

	r.initTemplates()
	return r
}

//Start method starts the AskMe App
func (app *App) Start() error {
	mux := httprouter.New()

	for _, a := range app.rtr.Actions() {
		log.Printf("Method %s, Path %s\n", a.Method, a.Path)
		mux.Handle(a.Method, a.Path, a.Func)
	}

	mux.ServeFiles("/public/*filepath", app.fileProvider)

	fmt.Println("Hello!")
	fmt.Printf("Listening on port %d\n", app.config.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", app.config.Port), mux)
}

// NewApp Creates a new app server
func NewApp(
	config *Config,
	ctrl Router, fileProvider shared.FileProvider) *App {
	return &App{
		config:       config,
		fileProvider: fileProvider,
		rtr:          ctrl,
	}
}
