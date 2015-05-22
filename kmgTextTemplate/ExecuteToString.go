package kmgTextTemplate

import (
	"bytes"
	"text/template"
)

func ExecuteToString(tmpl *template.Template, data interface{}) (output string, err error) {
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, data)
	if err != nil {
		return
	}
	return buf.String(), nil
}

func MustRenderToByte(text string, data interface{}) (b []byte) {
	w := &bytes.Buffer{}
	tmpl, err := template.New("").Parse(text)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
	return w.Bytes()
}
