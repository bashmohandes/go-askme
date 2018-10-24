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
	user.AsksUsecase
	user.AnswersUsecase
	sstore framework.SessionStore
}

// NewHomeController returns a new controller
func NewHomeController(
	rtr framework.Router,
	rndr framework.Renderer,
	askUC user.AsksUsecase,
	answrUC user.AnswersUsecase,
	sstr framework.SessionStore) *HomeController {
	c := &HomeController{
		Router:         rtr,
		Renderer:       rndr,
		AsksUsecase:    askUC,
		AnswersUsecase: answrUC,
		sstore:         sstr,
	}
	c.Get("/", c.index)
	c.Get("/me/top", c.topAnswers)
	c.Get("/questions", c.questions)
	c.Post("/questions", c.postQuestion)
	return c
}

// Index serves homepage
func (c *HomeController) index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c.Render(w, framework.ViewModel{BodyTmpl: "index", Title: "Index", Bag: framework.Map{}})
}

// TopAnswers serves top answers
func (c *HomeController) questions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c.Render(w, framework.ViewModel{BodyTmpl: "questions", Title: "Questions", Bag: framework.Map{}})
}

// TopAnswers serves top answers
func (c *HomeController) topAnswers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c.Render(w, framework.ViewModel{BodyTmpl: "top", Title: "Top Answers", Bag: framework.Map{}})
}

// PostQuestion posts a new question
func (c *HomeController) postQuestion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user1, _ := models.NewUser("test@test.com", "Test User", "")
	user2, _ := models.NewUser("bashmohandes@live.com", "Mohamed Elsherif", "")
	c.Ask(user1, user2, r.PostFormValue("question"))
	c.index(w, r, ps)
}
