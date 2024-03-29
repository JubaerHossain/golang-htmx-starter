package handler

import (
	"fmt"
	"net/http"

	"github.com/JubaerHossain/golang-htmx-starter/pkg/core"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// RootPath handles the root route to display index.gohtml
func RootPath(c echo.Context, a *core.App) error {
	// Retrieve the session
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}

	// Retrieve the value from the session
	helloValue := sess.Values["hello"]
	BuildVersion := sess.Values["build_version"]

	// Create a data map to pass to the renderer
	data := map[string]interface{}{
		"Title":        "Welcome to Htmx",
		"HelloSession": helloValue,
		"BuildVersion": BuildVersion,
	}

	fmt.Println("HelloSession: ", helloValue)
	fmt.Println("BuildVersion: ", BuildVersion)

	// Use Echo's built-in rendering method
	return c.Render(http.StatusOK, "index.html", data)
}
