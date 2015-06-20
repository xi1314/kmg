package kmgHTMLPurifier

import (
	"github.com/microcosm-cc/bluemonday"
)

// 去掉html属性
func HTMLPurifier(htmlInput string) string {
	p := bluemonday.UGCPolicy()
	html := p.Sanitize(htmlInput)
	return html
}

// 去掉所有html
func HTMLPurifierAll(htmlInput string) string {
	p := bluemonday.NewPolicy()
	html := p.Sanitize(htmlInput)
	return html
}
