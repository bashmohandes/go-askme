package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bashmohandes/go-askme/models"
	"github.com/bashmohandes/go-askme/user/usecase"
	"github.com/bashmohandes/go-askme/web/framework"
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
	c.Get("/u/:email", c.userFeed).Authenticated()
	c.Get("/u/:email/questions", c.questions).Authenticated()
	c.Post("/u/:email/questions", c.postQuestion).Authenticated()
	c.Post("/u/:email/answer/:questionId", c.postAnswer).Authenticated()
	c.Post("/u/:email/follow", c.follow).Authenticated()
	return c
}

// Me serves profile page
func (c *ProfileController) userFeed(cxt framework.Context) {
	email := cxt.Params().ByName("email")
	if len(email) == 0 {
		// flash message
		cxt.Redirect("/", http.StatusTemporaryRedirect)
	}
	profileUser, err := c.FindUserByEmail(email)
	if err != nil {
		// flash message
	}
	feed, err := c.LoadUserFeed(profileUser)
	if err != nil {
		cxt.ResponseWriter().Write([]byte(err.Error()))
		return
	}
	c.Render(
		cxt.ResponseWriter(),
		framework.ViewModel{
			BodyTmpl: "feed.body",
			Title:    "Home",
			Bag: framework.Map{
				"User":        cxt.Session().Get("user"),
				"ProfileUser": profileUser,
				"Feed":        feed}})
}

// TopAnswers serves top answers
func (c *ProfileController) questions(cxt framework.Context) {
	email := cxt.Params().ByName("email")
	if len(email) == 0 {
		// flash message
		cxt.Redirect("/", http.StatusTemporaryRedirect)
	}
	profileUser, err := c.FindUserByEmail(email)
	if err != nil {
		// flash message
	}
	feed, err := c.FetchUnansweredQuestions(profileUser)
	if err != nil {
		// flash message
	}
	c.Render(
		cxt.ResponseWriter(),
		framework.ViewModel{
			BodyTmpl: "profile.body",
			Title:    "Profile",
			Bag: framework.Map{
				"User":        cxt.Session().Get("user"),
				"ProfileUser": profileUser,
				"Feed":        feed}})
}

// PostQuestion posts a new question
func (c *ProfileController) postQuestion(cxt framework.Context) {
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
	user2, err := c.FindUserByEmail(email)
	if err != nil {
		// flash message
		cxt.Redirect("/", http.StatusTemporaryRedirect)
	}
	c.Ask(user1, user2, cxt.Request().PostFormValue("question"))
	c.questions(cxt)
}

func (c *ProfileController) postAnswer(cxt framework.Context) {
	user, ok := cxt.Session().Get("user").(*models.User)
	if !ok {
		// flash message
		cxt.Redirect("/", http.StatusTemporaryRedirect)
	}
	questionID64, err := strconv.ParseUint(cxt.Params().ByName("questionId"), 10, 32)
	if err != nil {
		// parse error
	}
	questionID := uint(questionID64)
	question, err := c.FetchQuestionById(questionID)
	if err != nil {
		// flash message
		cxt.Redirect("/", http.StatusTemporaryRedirect)
	}
	c.Answer(user, question, cxt.Request().PostFormValue("answer"))

	email := cxt.Params().ByName("email")
	redir := fmt.Sprintf("/u/%s", email)
	cxt.Redirect(redir, http.StatusFound)
}

func (c *ProfileController) follow(cxt framework.Context) {
	user, ok := cxt.Session().Get("user").(*models.User)
	if !ok {
		cxt.Redirect("/", http.StatusTemporaryRedirect)
	}

	toFollowEmail := cxt.Request().PostFormValue("toFollowEmail")
	if toFollowEmail == "" {
		// parse error
	}
	toFollow, err := c.FindUserByEmail(toFollowEmail)
	if err != nil {
		// flash error
	}

	err = c.Follow(user, toFollow)
	if err != nil {
		// flash error
	}

	cxt.Redirect("/", http.StatusSeeOther)
}
