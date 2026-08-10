package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bashmohandes/go-askme/models"
	"github.com/bashmohandes/go-askme/user/usecase"
	"github.com/bashmohandes/go-askme/web/framework"
	"github.com/bashmohandes/go-askme/web/middleware/csrf"
)

// ProfileController type
type ProfileController struct {
	framework.Router
	framework.Renderer
	user.AnswersUsecase
	user.AsksUsecase
}

// NewProfileController creates a new ProfileController
func NewProfileController(
	rtr framework.Router,
	rndr framework.Renderer,
	askUC user.AsksUsecase,
	answrUC user.AnswersUsecase) *ProfileController {
	c := &ProfileController{
		Router:         rtr,
		Renderer:       rndr,
		AsksUsecase:    askUC,
		AnswersUsecase: answrUC,
	}
	c.Get("/u/:email", c.userFeed)
	c.Get("/u/:email/questions", c.questions).Authenticated()
	c.Post("/u/:email/questions", framework.AdaptRoute(c.postQuestion, csrf.RequireCSRF)).Authenticated()
	c.Post("/u/:email/answer/:questionId", framework.AdaptRoute(c.postAnswer, csrf.RequireCSRF)).Authenticated()
	return c
}

// Me serves profile page
func (c *ProfileController) userFeed(cxt framework.Context) {
	email := cxt.Params().ByName("email")
	if len(email) == 0 {
		http.Error(cxt.ResponseWriter(), "profile not found", http.StatusNotFound)
		return
	}
	profileUser, err := c.FindUserByEmail(email)
	if err != nil {
		http.Error(cxt.ResponseWriter(), "profile not found", http.StatusNotFound)
		return
	}
	feed, err := c.LoadUserFeed(profileUser)
	if err != nil {
		http.Error(cxt.ResponseWriter(), "could not load profile", http.StatusInternalServerError)
		return
	}
	c.Render(
		cxt.ResponseWriter(),
		framework.ViewModel{
			BodyTmpl: "feed.body",
			Title:    "Home",
			Bag: framework.Map{
				"User":           cxt.Session().Get("user"),
				"ProfileUser":    profileUser,
				"Feed":           feed,
				csrf.TemplateTag: csrf.TemplateField(cxt)}})
}

// TopAnswers serves top answers
func (c *ProfileController) questions(cxt framework.Context) {
	email := cxt.Params().ByName("email")
	if len(email) == 0 {
		http.Error(cxt.ResponseWriter(), "profile not found", http.StatusNotFound)
		return
	}
	profileUser, err := c.FindUserByEmail(email)
	if err != nil {
		http.Error(cxt.ResponseWriter(), "profile not found", http.StatusNotFound)
		return
	}
	currentUser, ok := cxt.Session().Get("user").(*models.User)
	if !ok || currentUser.ID != profileUser.ID {
		http.Error(cxt.ResponseWriter(), "forbidden", http.StatusForbidden)
		return
	}
	feed, err := c.FetchUnansweredQuestions(profileUser)
	if err != nil {
		http.Error(cxt.ResponseWriter(), "could not load questions", http.StatusInternalServerError)
		return
	}
	c.Render(
		cxt.ResponseWriter(),
		framework.ViewModel{
			BodyTmpl: "profile.body",
			Title:    "Profile",
			Bag: framework.Map{
				"User":           cxt.Session().Get("user"),
				"ProfileUser":    profileUser,
				"Feed":           feed,
				csrf.TemplateTag: csrf.TemplateField(cxt)}})
}

// PostQuestion posts a new question
func (c *ProfileController) postQuestion(cxt framework.Context) {
	user1, ok := cxt.Session().Get("user").(*models.User)
	if !ok {
		http.Error(cxt.ResponseWriter(), "forbidden", http.StatusForbidden)
		return
	}
	email := cxt.Params().ByName("email")
	if len(email) == 0 {
		http.Error(cxt.ResponseWriter(), "profile not found", http.StatusNotFound)
		return
	}
	user2, err := c.FindUserByEmail(email)
	if err != nil {
		http.Error(cxt.ResponseWriter(), "profile not found", http.StatusNotFound)
		return
	}
	c.Ask(user1, user2, cxt.Request().PostFormValue("question"))
	cxt.Redirect(fmt.Sprintf("/u/%s", email), http.StatusSeeOther)
}

func (c *ProfileController) postAnswer(cxt framework.Context) {
	signedInUser, ok := cxt.Session().Get("user").(*models.User)
	if !ok {
		http.Error(cxt.ResponseWriter(), "forbidden", http.StatusForbidden)
		return
	}
	questionId64, err := strconv.ParseUint(cxt.Params().ByName("questionId"), 10, 32)
	if err != nil {
		http.Error(cxt.ResponseWriter(), "question not found", http.StatusNotFound)
		return
	}
	questionId := uint(questionId64)
	question, err := c.FetchQuestionById(questionId)
	if err != nil {
		http.Error(cxt.ResponseWriter(), "question not found", http.StatusNotFound)
		return
	}
	if _, err := c.Answer(signedInUser, question, cxt.Request().PostFormValue("answer")); err != nil {
		if errors.Is(err, user.ErrNotQuestionRecipient) {
			http.Error(cxt.ResponseWriter(), "forbidden", http.StatusForbidden)
			return
		}
		http.Error(cxt.ResponseWriter(), "could not save answer", http.StatusInternalServerError)
		return
	}

	email := cxt.Params().ByName("email")
	redir := fmt.Sprintf("/u/%s", email)
	cxt.Redirect(redir, http.StatusFound)
}
