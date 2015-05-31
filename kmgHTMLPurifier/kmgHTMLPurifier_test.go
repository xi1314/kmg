package kmgHTMLPurifier_test

import (
	"github.com/bronze1man/kmg/kmgHTMLPurifier"
	"github.com/bronze1man/kmg/kmgTest"
	"strings"
	"testing"
)

// 过滤 HTML的内容
func TestHTMLPurifier(t *testing.T) {
	input := []string{
		`<a onblur="alert(secret)" href="http://www.google.com">Google</a>`,
		`<a href="http://www.google.com/" target="_blank"><img src="https://ssl.gstatic.com/accounts/ui/logo_2x.png"/></a>`,
	}
	for _, content := range input {
		output := kmgHTMLPurifier.HTMLPurifier(content)
		kmgTest.Ok(!strings.Contains(output, "alert"))
	}
}
