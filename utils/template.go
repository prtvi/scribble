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
	dataMap := data.(map[string]any)

	if !IsProdEnv() {
		dataMap["debug"] = true
		dataMap["FF"] = FontFaces
	}

	if name != "app" {
		dataMap["AboutText"] = AboutText
		dataMap["HowToSlides"] = HowToSlides
	}

	return t.templates[name].ExecuteTemplate(w, name, data)
}

// render nested templates
func InitTemplates() *Template {
	tmpls := make(map[string]*template.Template)

	tmpls["app"] = template.Must(template.ParseFiles(
		"public/views/app.html",
		"public/views/partials/header.html",
	))

	tmpls["index"] = template.Must(template.ParseFiles(
		"public/views/index.html",
		"public/views/partials/header.html",
		"public/views/partials/footer.html",
		"public/views/partials/customiseAvatar.html",
	))

	return &Template{
		templates: tmpls,
	}
}
