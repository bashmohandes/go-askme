package controllers

import (
	"net/http"

	"github.com/bashmohandes/go-askme/web/framework"
	"github.com/julienschmidt/httprouter"
)

// ProfileController type
type ProfileController struct {
	framework.Router
	framework.Renderer
}

// NewProfileController creates a new ProfileController
func NewProfileController(rtr framework.Router, rndr framework.Renderer) *ProfileController {
	c := &ProfileController{
		Router:   rtr,
		Renderer: rndr,
	}
	c.Get("/me", c.index)

	return c
}

// Me serves profile page
func (c *ProfileController) index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c.Render(w, framework.ViewModel{BodyTmpl: "me", Title: "Me", Data: nil})
}
