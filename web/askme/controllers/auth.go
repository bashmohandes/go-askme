package controllers

import (
	"fmt"
	"net/http"

	"github.com/bashmohandes/go-askme/user/usecase"
	"github.com/bashmohandes/go-askme/web/framework"
	"github.com/bashmohandes/go-askme/web/middleware/csrf"
)

// AuthController manages authentication actions
type AuthController struct {
	framework.Router
	framework.Renderer
	user.AuthUsecase
	smgr framework.SessionManager
}

// NewAuthController creates a new AuthController
func NewAuthController(
	rtr framework.Router,
	rndr framework.Renderer,
	smgr framework.SessionManager,
	authUC user.AuthUsecase) *AuthController {
	c := &AuthController{
		Router:      rtr,
		Renderer:    rndr,
		smgr:        smgr,
		AuthUsecase: authUC,
	}

	c.Get("/login", c.login)
	c.Post("/login",
		framework.AdaptRoute(c.performLogin, csrf.RequireCSRF))
	c.Get("/signup", c.signup)
	c.Post("/signup",
		framework.AdaptRoute(c.performSignup, csrf.RequireCSRF))
	c.Get("/logout", c.logout).Authenticated()

	return c
}

func (c *AuthController) login(cxt framework.Context) {
	cxt.Session().Set("redir", cxt.Request().URL.Query().Get("redir"))
	c.Render(
		cxt.ResponseWriter(),
		framework.ViewModel{
			BodyTmpl: "standardlogin.body",
			Title: "Login", HeadTmpl: "standardlogin.head",
			Bag: framework.Map{
				csrf.TemplateTag: csrf.TemplateField(cxt),
			}})
}

func (c *AuthController) performLogin(cxt framework.Context) {
	email := cxt.Request().PostFormValue("email")
	pwd := cxt.Request().PostFormValue("password")
	// remember := r.PostFormValue("rememberMe")
	user, err := c.Signin(email, pwd)
	if err != nil {
		cxt.ResponseWriter().Write([]byte(fmt.Sprintf("Err: %v", err)))
		return
	}

	cxt.Session().Set("user", user)
	cxt.SetUser(&framework.User{ID: string(user.ID), Name: user.Name})
	redir, _ := cxt.Session().Get("redir").(string)
	if len(redir) == 0 {
		redir = fmt.Sprintf("/u/%s", user.Email)
	}
	cxt.Redirect(redir, http.StatusFound)
}

func (c *AuthController) signup(cxt framework.Context) {
	c.Render(
		cxt.ResponseWriter(),
		framework.ViewModel{
			BodyTmpl: "signup.body",
			Title: "Signup",
			HeadTmpl: "signup.head",
			Bag: framework.Map{
				csrf.TemplateTag: csrf.TemplateField(cxt),
			}})
}

func (c *AuthController) performSignup(cxt framework.Context) {
	r := cxt.Request()
	w := cxt.ResponseWriter()
	email := r.PostFormValue("email")
	pwd := r.PostFormValue("password")
	name := r.PostFormValue("name")
	_, err := c.Signup(email, pwd, name)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Err: %v", err)))
		return
	}
	cxt.Redirect("/login", http.StatusFound)
}

func (c *AuthController) logout(cxt framework.Context) {
	c.smgr.Abandon(cxt)
	cxt.Redirect("/", http.StatusFound)
}
