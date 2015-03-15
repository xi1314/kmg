package kmgHtmlTemplate

import (
	"html/template"
	"io"
)

func MustRenderToWriter(w io.Writer, text string, data interface{}) {
	tmpl := template.Must(template.New("").Parse(text))
	err := tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
	return
}
