package kmgFile

import (
	. "github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestKmgFile(ot *testing.T) {
	err := WriteFile(".kmgFileTest", []byte(""))
	Equal(err, nil)
	MustDeleteFile(".kmgFileTest")
	MustDeleteFile(".kmgFileTest")
	MustDeleteFileOrDirectory(".kmgFileTest")
}
