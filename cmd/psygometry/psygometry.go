package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/joho/godotenv"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

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

		answers, err := ParsePsychometryAnswers(req.Form, fakeData)
		if err != nil {
			return err
		}

		summary, err := CalculateScoreSummary(fakeData, *answers)
		if err != nil {
			return err
		}
		log.Println("score summary = ")
		summaryJson, err := json.Marshal(summary)
		log.Println(string(summaryJson))

		return c.Render(http.StatusCreated, "scores", summary)
	})
	e.Logger.Fatal(e.Start(":1714"))
}
