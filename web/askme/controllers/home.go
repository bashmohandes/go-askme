package controllers

import (
	"net/http"

	"github.com/bashmohandes/go-askme/model"
	"github.com/bashmohandes/go-askme/user/usecase"
	"github.com/bashmohandes/go-askme/web/framework"
	"github.com/julienschmidt/httprouter"
)

// HomeController represents Controller in MVC  model
type HomeController struct {
	framework.Router
	framework.Renderer
	asksUserScenario    user.AsksUsecase
	answerUserScenarion user.AnswersUsecase
}

// NewHomeController returns a new controller
func NewHomeController(rtr framework.Router, rndr framework.Renderer, askUC user.AsksUsecase, answrUC user.AnswersUsecase) *HomeController {
	c := &HomeController{
		Router:              rtr,
		Renderer:            rndr,
		asksUserScenario:    askUC,
		answerUserScenarion: answrUC,
	}
	c.Get("/", c.index)
	c.Get("/me/top", c.topAnswers)
	c.Get("/questions", c.questions)
	c.Post("/questions", c.postQuestion)

	return c
}

// Index serves homepage
func (c *HomeController) index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	questions := c.asksUserScenario.LoadUserFeed(models.NewUser("Visitor@hotmail.com", "Visitor Visiting"))
	c.Render(w, framework.ViewModel{Template: "index", Title: "Index", Data: questions})
}

// TopAnswers serves top answers
func (c *HomeController) questions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	questions := c.answerUserScenarion.FetchUnansweredQuestions(models.NewUniqueID())
	c.Render(w, framework.ViewModel{Template: "questions", Title: "Questions", Data: questions})
}

// TopAnswers serves top answers
func (c *HomeController) topAnswers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c.Render(w, framework.ViewModel{Template: "top", Title: "Top Answers", Data: nil})
}

// PostQuestion posts a new question
func (c *HomeController) postQuestion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	c.asksUserScenario.Ask(models.NewUser("test@test.com", "Test User"), models.NewUser("bashmohandes@live.com", "Mohamed Elsherif"), r.FormValue("question"))
	c.index(w, r, ps)
}
