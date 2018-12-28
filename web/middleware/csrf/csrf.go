package csrf

import (
	"fmt"
	"net/http"
	"html/template"

	"github.com/bashmohandes/go-askme/web/framework"
	"github.com/google/uuid"
)

const (
	fieldName	string = "csrf_token"
	sessionKey  string = "_csrf_token" 
)

var TemplateTag = "csrfField"

// CSRF adds a csrf token to all sessions
type csrf struct {
}

// Add CSRF token if it does not exist in session
func (c *csrf) Run(cxt framework.Context) bool {
	session := cxt.Session()
	if session.Get(sessionKey) == nil {
		session.Set(sessionKey, uuid.New().String())
	}
	return true
}

func CSRF() framework.Middleware {
	return &csrf{}
}

// Decorator for framework.RouteHandler that checks the
// validity of the CSRF token by comparing the token in the
// form (hidden input) and the correct token in the session.
// It gives 403 Forbidden error if it is invalid and calls the
// given RouteHandler otherwise.
// E.g., router.Post("/login", csrf.RequireCSRF(performLogin))
// where performLogin is of type framework.RouteHandler
func RequireCSRF(f framework.RouteHandler) framework.RouteHandler {
	return func(cxt framework.Context) {
		sessionCSRF := cxt.Session().Get(sessionKey)
		formCSRF := cxt.Request().PostFormValue(fieldName)
		if sessionCSRF != formCSRF {
			w := cxt.ResponseWriter()
			http.Error(w, "403 - Forbidden!", http.StatusForbidden)
			return
		}
		f(cxt)
	}
}

func TemplateField(cxt framework.Context) template.HTML {
	sessionCSRF := cxt.Session().Get(sessionKey)
	if sessionCSRF != nil {
		fragment := fmt.Sprintf(`<input type="hidden" name="%s" value="%s">`,
			fieldName, sessionCSRF)

		return template.HTML(fragment)
	}

	return template.HTML("")
}