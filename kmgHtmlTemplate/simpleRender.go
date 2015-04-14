package kmgHtmlTemplate

import (
	"bytes"
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

func RenderToWriter(w io.Writer, text string, data interface{}) (err error) {
	tmpl, err := template.New("").Parse(text)
	if err != nil {
		return
	}
	return tmpl.Execute(w, data)
}

func RenderToByte(text string, data interface{}) (b []byte, err error) {
	w := &bytes.Buffer{}
	err = RenderToWriter(w, text, data)
	if err != nil {
		return
	}
	return w.Bytes(), nil
}
