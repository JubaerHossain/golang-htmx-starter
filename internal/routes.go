package serve

import (
	"github.com/JubaerHossain/golang-htmx-starter/internal/handler"
	"github.com/JubaerHossain/golang-htmx-starter/pkg/core"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// bindRootRoute binds the root route
func BindRootRoute(a *core.App, path string) {
	a.Echo.GET(path, func(c echo.Context) error {
		// set hello world in session
		sess, _ := session.Get("session", c)
		sess.Values["hello"] = "world"
		sess.Values["build_version"] = "v0.0.1"
		sess.Save(c.Request(), c.Response())
		return handler.RootPath(c, a)
	})
}
