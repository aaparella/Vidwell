package render

import (
	"encoding/json"
	"html/template"
	"io"
	"path"

	"github.com/aaparella/vidwell/config"
)

var templatesDir string

func init() {
	templatesDir = config.GetRenderingConfiguration().TemplatesDir
	// TODO: add verification that this directory exists, etc.
}

// Render writes the template, rendered with data, to the passwed writer.
// Renders an error page if the template specified does not exist (should be
// the name of a template file in views/templates/, without the .tpl),
// or cannot be opened. Will then return any errors caused by parsing
// the template, or rendering it with the provided data.
func Render(w io.Writer, tmpl string, data interface{}) {
	template, err := template.ParseFiles(path.Join(templatesDir, tmpl+".tmpl"))
	if err != nil {
		renderErrorPage(w, tmpl, data, err)
		return
	}
	err = template.Execute(w, data)
	if err != nil {
		renderErrorPage(w, tmpl, data, err)
	}
}

// ErrorPageData is used to render an error page when an error occurs while
// rendering a page.
type ErrorPageData struct {
	File  string
	Data  string
	Error string
}

// renderErrorPage displays what the error was, along with the data passed
// in to the template, along with the template file name.
//
// This in turn calls Render, which may not be the best idea,
// because if *that* call fails, there is an endless recursion where Render
// and renderErrorPage keep getting called. So, make sure renderErrorPage
// cannot cause an error with it's call to Render.
func renderErrorPage(w io.Writer, tmpl string, data interface{}, err error) {
	d, _ := json.MarshalIndent(data, "", "	")
	Render(w, "error", ErrorPageData{
		File:  tmpl,
		Data:  string(d),
		Error: err.Error(),
	})
}
