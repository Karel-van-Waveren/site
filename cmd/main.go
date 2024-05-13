package main

import (
	"html/template"
	"io"

	// "log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type Template struct {
	tmpl *template.Template
}

func newTemplate() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("views/*.html")),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

type State struct{}

func NewData() *State {
	return &State{}
}

func init() {
}

func main() {
	debug := true
	err := godotenv.Load()
	if err == nil {
		debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	}

	e := echo.New()

	if debug {
		e.Logger.SetLevel(log.DEBUG)
	}

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	data := NewData()

	e.Renderer = newTemplate()
	e.Use(middleware.Logger())
	e.Static("/css", "css")

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index.html", data)
	})

	if debug {
		e.Logger.Fatal(e.Start(":80"))
	} else {
		fullchain := os.Getenv("FULLCHAIN_PATH")
		privkey := os.Getenv("PRIVKEY_PATH")
		log.Print(fullchain)
		log.Print(privkey)

		err := e.StartTLS(":433", fullchain, privkey)
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}
}
