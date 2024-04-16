package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	e.Renderer = t

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", generateFakeData())
	})
	e.POST("/answers", func(c echo.Context) error {
		req := c.Request()
		err := req.ParseForm()
		if err != nil {
			return err
		}

		fmt.Print("----------\n")
		for x := range req.Form {
			fmt.Printf("'%s' = '%s'\n", x, req.Form.Get(x))
		}
		fmt.Print("----------\n")

		return c.NoContent(http.StatusCreated)
	})
	e.Logger.Fatal(e.Start(":1714"))
}
