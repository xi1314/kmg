package JsTpl

import (
	"testing"

	"github.com/bronze1man/kmg/kmgTest"
)

func TestMustBuildTplOneFile(ot *testing.T) {
	kmgTest.Equal(MustBuildTplOneFile([]byte(``)), nil)
	kmgTest.Equal(MustBuildTplOneFile([]byte("var a = ``;")), []byte(`var a = "";`))
	kmgTest.Equal(MustBuildTplOneFile([]byte("var a = `a\na`;")), []byte(`var a = "a\na";`))
}
