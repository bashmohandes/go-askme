package controllers

import (
	"fmt"
	"net/http"

	"github.com/bashmohandes/go-askme/user/usecase"
	"github.com/bashmohandes/go-askme/web/framework"
	"github.com/julienschmidt/httprouter"
)

// AuthController manages authentication actions
type AuthController struct {
	framework.Router
	framework.Renderer
	user.AuthUsecase
	sstore framework.SessionStore
}

// NewAuthController creates a new AuthController
func NewAuthController(
	rtr framework.Router,
	rndr framework.Renderer,
	sstr framework.SessionStore,
	authUC user.AuthUsecase) *AuthController {
	c := &AuthController{
		Router:      rtr,
		Renderer:    rndr,
		sstore:      sstr,
		AuthUsecase: authUC,
	}

	c.Get("/login", c.login)
	c.Post("/login", c.performLogin)
	c.Get("/signup", c.signup)
	c.Post("/signup", c.performSignup)

	return c
}

func (c *AuthController) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c.Render(w, framework.ViewModel{BodyTmpl: "login.body", Title: "Login", HeadTmpl: "login.head", Bag: framework.Map{}})
}

func (c *AuthController) performLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	email := r.PostFormValue("email")
	pwd := r.PostFormValue("password")
	// remember := r.PostFormValue("rememberMe")
	user, err := c.Signin(email, pwd)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Err: %v", err)))
		return
	}
	session := c.sstore.FetchOrCreate(w, r)
	session.Set("user", user)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (c *AuthController) signup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c.Render(w, framework.ViewModel{BodyTmpl: "signup.body", Title: "Signup", HeadTmpl: "signup.head", Bag: framework.Map{}})
}

func (c *AuthController) performSignup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	email := r.PostFormValue("email")
	pwd := r.PostFormValue("password")
	name := r.PostFormValue("name")
	_, err := c.Signup(email, pwd, name)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Err: %v", err)))
		return
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}
