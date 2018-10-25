package controllers

import (
	"github.com/bashmohandes/go-askme/model"
	"github.com/bashmohandes/go-askme/user/usecase"
	"github.com/bashmohandes/go-askme/web/framework"
)

// HomeController represents Controller in MVC  model
type HomeController struct {
	framework.Router
	framework.Renderer
	user.AsksUsecase
	user.AnswersUsecase
	sstore framework.SessionManager
}

// NewHomeController returns a new controller
func NewHomeController(
	rtr framework.Router,
	rndr framework.Renderer,
	askUC user.AsksUsecase,
	answrUC user.AnswersUsecase) *HomeController {
	c := &HomeController{
		Router:         rtr,
		Renderer:       rndr,
		AsksUsecase:    askUC,
		AnswersUsecase: answrUC,
	}
	c.Get("/", c.index).Authenticated()
	c.Get("/me/top", c.topAnswers).Authenticated()
	c.Get("/questions", c.questions).Authenticated()
	c.Post("/questions", c.postQuestion).Authenticated()
	return c
}

// Index serves homepage
func (c *HomeController) index(cxt framework.Context) {
	c.Render(cxt.ResponseWriter(), framework.ViewModel{BodyTmpl: "index", Title: "Index", Bag: framework.Map{"user": cxt.Session().Get("user")}})
}

// TopAnswers serves top answers
func (c *HomeController) questions(cxt framework.Context) {
	c.Render(cxt.ResponseWriter(), framework.ViewModel{BodyTmpl: "questions", Title: "Questions", Bag: framework.Map{"user": cxt.Session().Get("user")}})
}

// TopAnswers serves top answers
func (c *HomeController) topAnswers(cxt framework.Context) {
	c.Render(cxt.ResponseWriter(), framework.ViewModel{BodyTmpl: "top", Title: "Top Answers", Bag: framework.Map{"user": cxt.Session().Get("user")}})
}

// PostQuestion posts a new question
func (c *HomeController) postQuestion(cxt framework.Context) {
	user1, _ := models.NewUser("test@test.com", "Test User", "")
	user2, _ := models.NewUser("bashmohandes@live.com", "Mohamed Elsherif", "")
	c.Ask(user1, user2, cxt.Request().PostFormValue("question"))
	c.index(cxt)
}
