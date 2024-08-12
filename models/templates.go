package models

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Templates struct {
	Templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() *Templates {
	return &Templates{
		Templates: template.Must(template.ParseGlob("templates/*.html")),
	}
}
