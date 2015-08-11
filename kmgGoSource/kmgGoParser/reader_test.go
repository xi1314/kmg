package kmgGoParser

import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestGoSourceRemoveComment(ot *testing.T) {
	content := []byte("// abc\n/*  abc\n*/\npackage abc\n\nfunc b(){\n}")
	outContent := goSourceRemoveComment(content, nil)
	kmgTest.Equal(string(outContent), "      \n       \n  \npackage abc\n\nfunc b(){\n}")
}
