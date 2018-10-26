package controllers

import (
	"net/http"

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
	c.Post("/u/:email/questions", c.postQuestion).Authenticated()
	return c
}

// Index serves homepage
func (c *HomeController) index(cxt framework.Context) {
	c.Render(cxt.ResponseWriter(), framework.ViewModel{BodyTmpl: "index", Title: "Index", Bag: framework.Map{"user": cxt.Session().Get("user")}})
}

// TopAnswers serves top answers
func (c *HomeController) questions(cxt framework.Context) {
	feed, err := c.FetchUnansweredQuestions(cxt.User().ID)
	if err != nil {
		// flash message
	}
	c.Render(cxt.ResponseWriter(), framework.ViewModel{BodyTmpl: "questions", Title: "Questions", Bag: framework.Map{"user": cxt.Session().Get("user"), "feed": feed}})
}

// TopAnswers serves top answers
func (c *HomeController) topAnswers(cxt framework.Context) {
	c.Render(cxt.ResponseWriter(), framework.ViewModel{BodyTmpl: "top", Title: "Top Answers", Bag: framework.Map{"user": cxt.Session().Get("user")}})
}

// PostQuestion posts a new question
func (c *HomeController) postQuestion(cxt framework.Context) {
	user1, ok := cxt.Session().Get("user").(*models.User)
	if !ok {
		// flash message
		cxt.Redirect("/", http.StatusTemporaryRedirect)
	}
	email := cxt.Params().ByName("email")
	if len(email) == 0 {
		// flash message
		cxt.Redirect("/", http.StatusTemporaryRedirect)
	}
	user2 := c.FindUserByEmail(email)
	if user2 == nil {
		// flash message
		cxt.Redirect("/", http.StatusTemporaryRedirect)
	}
	c.Ask(user1, user2, cxt.Request().PostFormValue("question"))
	c.index(cxt)
}
