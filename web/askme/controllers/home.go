package controllers

import (
	"fmt"
	"net/http"

	"github.com/bashmohandes/go-askme/models"
	"github.com/bashmohandes/go-askme/web/framework"
)

// HomeController represents Controller in MVC  model
type HomeController struct {
	framework.Router
}

// NewHomeController returns a new controller
func NewHomeController(rtr framework.Router) *HomeController {
	c := &HomeController{rtr}
	c.Get("/", c.index).Authenticated()
	return c
}

// Index serves homepage
func (c *HomeController) index(cxt framework.Context) {
	user := cxt.Session().Get("user").(*models.User)
	http.Redirect(cxt.ResponseWriter(), cxt.Request(), fmt.Sprintf("/u/%s", user.Email), http.StatusFound)
}
