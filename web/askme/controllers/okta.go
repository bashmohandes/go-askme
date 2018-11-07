package controllers

import (
	"net/http"

	"github.com/bashmohandes/go-askme/web/framework"
	"github.com/bashmohandes/go-askme/web/oktautils"
)

// OktaController handles OKTA login functionality
type OktaController struct {
	framework.Router
	framework.Renderer
	config *framework.Config
}

var (
	state = "ApplicationState"
	nonce = "NonceNotSetYet"
)

// NewOktaController creates OKTA Controller
func NewOktaController(rtr framework.Router, rndr framework.Renderer, config *framework.Config) *OktaController {
	c := &OktaController{
		Router:   rtr,
		Renderer: rndr,
		config:   config,
	}

	c.Get("/login", c.login)
	c.Get("/authorization-code/callback", c.callback)
	c.Get("/logout", c.logout)
	return c
}

func (o *OktaController) login(cxt framework.Context) {
	r := cxt.Request()
	w := cxt.ResponseWriter()
	nonce, _ = oktautils.GenerateNonce()
	var redirectPath string

	q := r.URL.Query()
	q.Add("client_id", o.config.OktaClient)
	q.Add("response_type", "code")
	q.Add("response_mode", "query")
	q.Add("scope", "openid profile email")
	q.Add("redirect_uri", "http://localhost:8080/authorization-code/callback")
	q.Add("state", state)
	q.Add("nonce", nonce)

	redirectPath = o.config.OktaIssuer + "/v1/authorize?" + q.Encode()

	http.Redirect(w, r, redirectPath, http.StatusMovedPermanently)
}

func (o *OktaController) callback(cxt framework.Context) {

}

func (o *OktaController) logout(cxt framework.Context) {

}
