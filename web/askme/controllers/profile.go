package controllers

import (
	"github.com/bashmohandes/go-askme/web/framework"
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
	c.Get("/me", c.index).Authenticated()

	return c
}

// Me serves profile page
func (c *ProfileController) index(cxt framework.Context) {
	w := cxt.ResponseWriter()
	c.Render(w, framework.ViewModel{BodyTmpl: "questions", Title: "Me", Bag: framework.Map{}})
}
