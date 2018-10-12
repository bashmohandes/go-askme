package controllers

import (
	"fmt"
	"net/http"

	"github.com/bashmohandes/go-askme/model"
	"github.com/bashmohandes/go-askme/shared"
	"github.com/bashmohandes/go-askme/user/usecase"
	"github.com/bashmohandes/go-askme/web"
	"github.com/julienschmidt/httprouter"
)

// HomeController represents Controller in MVC  model
type HomeController struct {
	framework.Controller
	asksUserScenario    user.AsksUsecase
	answerUserScenarion user.AnswersUsecase
}

// NewHomeController returns a new controller
func NewHomeController(askUC user.AsksUsecase, answrUC user.AnswersUsecase, fp shared.FileProvider, config *framework.Config) *HomeController {
	c := &HomeController{
		asksUserScenario:    askUC,
		answerUserScenarion: answrUC,
	}
	c.Init(fp, config)
	c.AddAction("GET", "/", c.index)
	c.AddAction("GET", "/me/top", c.topAnswers)
	c.AddAction("POST", "/question", c.postQuestion)

	return c
}

// Index serves homepage
func (c *HomeController) index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	questions := c.asksUserScenario.LoadUserFeed(models.NewUser("Visitor@hotmail.com", "Visitor Visiting"))
	c.Render(w, framework.ViewModel{Template: "index", Title: "Index", Data: questions})
}

// TopAnswers serves top answers
func (c *HomeController) topAnswers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c.Render(w, framework.ViewModel{Template: "top", Title: "Top Answers", Data: nil})
}

// PostQuestion posts a new question
func (c *HomeController) postQuestion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()

	fmt.Fprintf(w, "Params %v", r.FormValue("username"))
}
