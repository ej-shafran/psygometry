package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

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

	fakeData := generateFakeData()

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", fakeData)
	})
	e.POST("/answers", func(c echo.Context) error {
		req := c.Request()
		err := req.ParseForm()
		if err != nil {
			return err
		}

		fmt.Print("----------\n")
		answers, err := ParsePsychometryAnswers(req.Form, fakeData)
		if err != nil {
			return err
		}
		json.NewEncoder(os.Stdout).Encode(answers)
		fmt.Print("----------\n")

		return c.NoContent(http.StatusCreated)
	})
	e.Logger.Fatal(e.Start(":1714"))
}
