package utils

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates map[string]*template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates[name].ExecuteTemplate(w, name, data)
}

var T *Template

// render nested templates
func InitTemplates() *Template {
	tmpls := make(map[string]*template.Template)

	tmpls["app"] = template.Must(template.ParseFiles(
		"public/views/app.html",
		"public/views/partials/header.html",
	))

	tmpls["createRoom"] = template.Must(template.ParseFiles(
		"public/views/createRoom.html",
		"public/views/partials/header.html",
	))

	tmpls["welcome"] = template.Must(template.ParseFiles(
		"public/views/welcome.html",
		"public/views/partials/header.html",
	))

	return &Template{
		templates: tmpls,
	}
}
