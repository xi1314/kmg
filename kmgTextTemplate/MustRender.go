package kmgTextTemplate

import (
	"bytes"
	"text/template"
)

//直接渲染一个模板(减少调用复杂度)
func MustRender(text string, data interface{}) []byte {
	tmpl := template.Must(template.New("").Parse(text))
	buf := &bytes.Buffer{}
	err := tmpl.Execute(buf, data)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
