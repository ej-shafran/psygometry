package main

import (
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/google/uuid"
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

type State struct {
	Page        int
	Session     string
	Psychometry Psychometry
}

var psychometries = map[string]*State{}
var formValues = map[string]*url.Values{}

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
		req := c.Request()
		if err := req.ParseForm(); err != nil {
			return err
		}

		session := req.Form.Get("session")
		if session == "" {
			session = uuid.New().String()
		}
		state, ok := psychometries[session]
		if !ok || state.Page >= len(state.Psychometry.Sections) {
			psychometry := generateFakeData()

			state = &State{
				Session:     session,
				Page:        -1,
				Psychometry: psychometry,
			}
			psychometries[session] = state
		}

		if state.Page < 0 {
			return c.Render(http.StatusOK, "writing-page", state)
		} else {
			return c.Render(http.StatusOK, "section-page", state)
		}
	})

	e.POST("/answers", func(c echo.Context) error {
		req := c.Request()
		if err := req.ParseForm(); err != nil {
			return err
		}

		session := req.Form.Get("session")
		if session == "" {
			session = uuid.New().String()
		}
		state, ok := psychometries[session]
		if !ok {
			// TODO: handle this
			return errors.New("invalid session")
		}

		values, ok := formValues[session]
		if !ok {
			values = &req.Form
			formValues[session] = values
		} else {
			for key, value := range req.Form {
				values.Set(key, value[0])
			}
		}

		state.Page += 1
		if state.Page < len(state.Psychometry.Sections) {
			return c.Render(http.StatusOK, "section", state.Psychometry.Sections[state.Page])
		}

		answers, err := ParsePsychometryAnswers(*values, state.Psychometry)
		if err != nil {
			return err
		}

		summary, err := CalculateScoreSummary(fakeData, *answers)
		if err != nil {
			return err
		}

		return c.Render(http.StatusCreated, "scores", summary)
	})

	e.Logger.Fatal(e.Start(":1714"))
}
