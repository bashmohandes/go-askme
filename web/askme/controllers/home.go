package controllers

import (
	"fmt"
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
	user.AuthUsecase
}

// NewHomeController returns a new controller
func NewHomeController(
	rtr framework.Router,
	rndr framework.Renderer,
	askUC user.AsksUsecase,
	answrUC user.AnswersUsecase,
	authUC user.AuthUsecase) *HomeController {
	c := &HomeController{
		Router:         rtr,
		Renderer:       rndr,
		AsksUsecase:    askUC,
		AnswersUsecase: answrUC,
		AuthUsecase:    authUC,
	}
	c.Get("/", c.index)
	c.Get("/me/top", c.topAnswers)
	c.Get("/questions", c.questions)
	c.Post("/questions", c.postQuestion)
	c.Get("/login", c.login)
	c.Post("/login", c.performLogin)
	c.Get("/signup", c.signup)
	c.Post("/signup", c.performSignup)
	return c
}

// Index serves homepage
func (c *HomeController) index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user, _ := models.NewUser("Visitor@hotmail.com", "Visitor Visiting", "")
	questions := c.LoadUserFeed(user)
	c.Render(w, framework.ViewModel{BodyTmpl: "index", Title: "Index", Data: questions})
}

// TopAnswers serves top answers
func (c *HomeController) questions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	questions := c.FetchUnansweredQuestions(models.NewUniqueID())
	c.Render(w, framework.ViewModel{BodyTmpl: "questions", Title: "Questions", Data: questions})
}

// TopAnswers serves top answers
func (c *HomeController) topAnswers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c.Render(w, framework.ViewModel{BodyTmpl: "top", Title: "Top Answers", Data: nil})
}

// PostQuestion posts a new question
func (c *HomeController) postQuestion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user1, _ := models.NewUser("test@test.com", "Test User", "")
	user2, _ := models.NewUser("bashmohandes@live.com", "Mohamed Elsherif", "")
	c.Ask(user1, user2, r.PostFormValue("question"))
	c.index(w, r, ps)
}

func (c *HomeController) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c.Render(w, framework.ViewModel{BodyTmpl: "login.body", Title: "Login", HeadTmpl: "login.head", Data: nil})
}

func (c *HomeController) performLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	email := r.PostFormValue("email")
	pwd := r.PostFormValue("password")
	// remember := r.PostFormValue("rememberMe")
	user, err := c.Signin(email, pwd)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Err: %v", err)))
		return
	}
	w.Write([]byte(fmt.Sprintf("Email: %s, Password: %s, Name: %s", user.Email, user.HashedPassword, user.Name)))
}

func (c *HomeController) signup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c.Render(w, framework.ViewModel{BodyTmpl: "signup.body", Title: "Signup", HeadTmpl: "signup.head", Data: nil})
}

func (c *HomeController) performSignup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	email := r.PostFormValue("email")
	pwd := r.PostFormValue("password")
	name := r.PostFormValue("name")
	user, err := c.Signup(email, pwd, name)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Err: %v", err)))
		return
	}
	w.Write([]byte(fmt.Sprintf("User: %v", user.ID)))
}
