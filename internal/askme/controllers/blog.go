package controllers

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/bashmohandes/go-askme/internal/askme/models"
)

type pageModel struct {
	Template string
	Title    string
	Data     interface{}
}

var tpl *template.Template

//Blog represents the main app model
func init() {
	tpl = template.New("master")
	tpl.Funcs(map[string]interface{}{
		"RenderTemplate": renderTemplate,
	})
	tpl.ParseGlob("../../internal/askme/templates/*")
}

func renderTemplate(name string, data interface{}) (ret template.HTML, err error) {
	buf := bytes.NewBuffer([]byte{})
	err = tpl.ExecuteTemplate(buf, name, data)
	ret = template.HTML(buf.String())
	return
}

func index(w http.ResponseWriter, r *http.Request) {
	questions := models.LoadQuestions("Bashmohandes")
	tpl.ExecuteTemplate(w, "master", pageModel{"index", "Index", questions})
}

func me(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "master", pageModel{"me", "Me", nil})
}

//Blog returns a new blog
func Blog() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/me", me)
	return mux
}
