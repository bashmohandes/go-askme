package controllers

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/bashmohandes/go-askme/internal/domain"
)

type pageModel struct {
	Template string
	Title    string
	Data     interface{}
}

var tpl *template.Template

//Blog represents the main app model
func init() {
	tpl = template.New("master")
	tpl.Funcs(map[string]interface{}{
		"RenderTemplate": renderTemplate,
	})
	tpl.ParseGlob("../../internal/askme/templates/*")
}

func renderTemplate(name string, data interface{}) (ret template.HTML, err error) {
	buf := bytes.NewBuffer([]byte{})
	err = tpl.ExecuteTemplate(buf, name, data)
	ret = template.HTML(buf.String())
	return
}

func index(w http.ResponseWriter, r *http.Request) {
	questions := models.LoadQuestions(models.NewUniqueID())
	render(w, pageModel{"index", "Index", questions})
}

func me(w http.ResponseWriter, r *http.Request) {
	render(w, pageModel{"me", "Me", nil})
}

func topUserAnswers(w http.ResponseWriter, r *http.Request) {
	render(w, pageModel{"top", "Top Answers", nil})
}

func render(w http.ResponseWriter, p pageModel) {
	tpl.ExecuteTemplate(w, "master", p)
}

//Blog returns a new blog
func Blog() http.Handler {
	mux := httprouter.New()
	mux.HandlerFunc("GET", "/", index)
	mux.HandlerFunc("GET", "/me", me)
	mux.HandlerFunc("GET", "/users/:userid/top", topUserAnswers)
	return mux
}
