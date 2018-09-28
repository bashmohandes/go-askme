package controllers

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/bashmohandes/go-askme/internal/domain"
	"github.com/bashmohandes/go-askme/internal/shared"
)

type pageModel struct {
	Template string
	Title    string
	Data     interface{}
}

var tpl *template.Template
var questionService QuestionService

//Blog represents the main app model
func init() {
	tpl = template.New("master")
	tpl.Funcs(map[string]interface{}{
		"RenderTemplate": renderTemplate,
	})
}

func renderTemplate(name string, data interface{}) (ret template.HTML, err error) {
	buf := bytes.NewBuffer([]byte{})
	err = tpl.ExecuteTemplate(buf, name, data)
	ret = template.HTML(buf.String())
	return
}

func index(w http.ResponseWriter, r *http.Request) {
	questions := questionService.LoadQuestions(models.NewUniqueID())
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

// QuestionService defines questions interface
type QuestionService interface {
	LoadQuestions(userID models.UniqueID) []models.Question
}

//Blog returns a new blog
func Blog(qs QuestionService, fp common.FileProvider) http.Handler {
	questionService = qs
	for _, t := range fp.List() {
		if strings.HasSuffix(t, ".gohtml") {
			tpl.Parse(fp.String(t))
		}
	}
	mux := httprouter.New()
	mux.HandlerFunc("GET", "/", index)
	mux.HandlerFunc("GET", "/me", me)
	mux.HandlerFunc("GET", "/users/:userid/top", topUserAnswers)
	return mux
}
