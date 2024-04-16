package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type State struct {
	Count int
}

func NewState() State {
	return State{
		Count: 1,
	}
}

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

	state := NewState()

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", state)
	})
	e.POST("/clicked", func(c echo.Context) error {
		state.Count += 1
		return c.Render(http.StatusCreated, "count", state)
	})
	e.Logger.Fatal(e.Start(":1714"))
}
