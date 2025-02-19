package utils

import (
	"html/template"
	"io"
	"io/fs"

	"github.com/labstack/echo/v4"
)

/*
=== TEMPLATE RENDERER ===
The html/template package is used to render HTML templates.
It provides a way to define templates and execute them to generate HTML output.

The echo.Renderer interface is used to render templates in Echo.
It has a single method, Render, which is used to render a template.

Here is the implementation of the echo.Renderer interface wrapped around the html/template.Template type.
=========================
*/

type TemplateRenderer struct {
	tmpl *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

func NewTemplateRendererFromFS(fs fs.FS, path string) echo.Renderer {
	return &TemplateRenderer{
		tmpl: template.Must(template.ParseFS(fs, path)),
	}
}
