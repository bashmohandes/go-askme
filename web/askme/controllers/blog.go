package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/bashmohandes/go-askme/model"
	"github.com/bashmohandes/go-askme/shared"
	"github.com/bashmohandes/go-askme/user/usecase"
	"github.com/julienschmidt/httprouter"
)

type pageModel struct {
	Template string
	Title    string
	Data     interface{}
}

var tpl *template.Template
var asksUserScenario user.AsksUsecase
var answerUserScenarion user.AnswersUsecase

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
	questions := asksUserScenario.LoadUserFeed(models.NewUser("Visitor@hotmail.com", "Visitor Visiting"))
	render(w, pageModel{"index", "Index", questions})
}

func me(w http.ResponseWriter, r *http.Request) {
	render(w, pageModel{"me", "Me", nil})
}

func topUserAnswers(w http.ResponseWriter, r *http.Request) {
	render(w, pageModel{"top", "Top Answers", nil})
}

func postQuestion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()

	fmt.Fprintf(w, "Params %v", r.FormValue("username"))
}

func render(w http.ResponseWriter, p pageModel) {
	err := tpl.ExecuteTemplate(w, "master", p)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

// QuestionService defines questions interface
type QuestionService interface {
	LoadQuestions(userID models.UniqueID) []*models.Question
}

//Blog returns a new blog
func Blog(askUC user.AsksUsecase, answrUC user.AnswersUsecase, fp shared.FileProvider) http.Handler {
	asksUserScenario = askUC
	answerUserScenarion = answrUC
	for _, t := range fp.List() {
		if strings.HasSuffix(t, ".gohtml") {
			tpl.Parse(fp.String(t))
		}
	}
	mux := httprouter.New()
	mux.HandlerFunc("GET", "/", index)
	mux.HandlerFunc("GET", "/me", me)
	mux.HandlerFunc("GET", "/me/top", topUserAnswers)
	mux.POST("/question/", postQuestion)
	mux.ServeFiles("/public/*filepath", fp)
	return mux
}
