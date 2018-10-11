package controllers

import (
	"net/http"

	"github.com/bashmohandes/go-askme/shared"

	"github.com/bashmohandes/go-askme/web"
	"github.com/julienschmidt/httprouter"
)

// ProfileController type
type ProfileController struct {
	framework.Controller
}

// NewProfileController creates a new ProfileController
func NewProfileController(fp shared.FileProvider) *ProfileController {
	c := &ProfileController{}
	c.Init(fp)
	c.AddAction("GET", "/me", c.index)

	return c
}

// Me serves profile page
func (c *ProfileController) index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c.Render(w, framework.ViewModel{Template: "me", Title: "Me", Data: nil})
}
