package kmgHTMLPurifier

import (
	"github.com/microcosm-cc/bluemonday"
)

func HTMLPurifier(htmlInput string) string {
	p := bluemonday.UGCPolicy()
	html := p.Sanitize(htmlInput)
	return html
}
