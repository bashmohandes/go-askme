package askme

import (
	"github.com/bashmohandes/go-askme/web/askme/controllers"
	"github.com/bashmohandes/go-askme/web/framework"
	"github.com/bashmohandes/go-askme/web/middleware/csrf"
)

// App represents the AskMe web app
type App struct {
	framework.App
}

// NewApp creates a new instance of Ask App
func NewApp(base framework.App,
	hc *controllers.HomeController,
	pc *controllers.ProfileController,
	ac *controllers.AuthController) *App {
	app := &App{
		App: base,
	}

	//app.App.Use(middleware.NewRequestLogger())
	app.App.Use(csrf.CSRF())
	return app
}
