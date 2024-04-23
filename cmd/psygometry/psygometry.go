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
		summary := CalculateScoreSummary(fakeData, *answers)

		log.Println("score summary = ")
		summaryJson, err := json.Marshal(summary)
		log.Println(string(summaryJson))

		essayScore, err := CalculateEssayScore(fakeData.EssaySection, answers.EssaySection)
		if err != nil {
			return err
		}

		log.Println("essay score = ")
		essayScoreJson, err := json.Marshal(essayScore)
		log.Println(string(essayScoreJson))

		return c.NoContent(http.StatusCreated)
	})
	e.Logger.Fatal(e.Start(":1714"))
}
