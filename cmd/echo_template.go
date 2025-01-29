package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Template struct {
	tmpl *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

// newTemplate return a echo.Renderer which is bound to the specified path
func newTemplate(path string) echo.Renderer {
	return &Template{
		tmpl: template.Must(template.ParseGlob(path)),
	}
}
