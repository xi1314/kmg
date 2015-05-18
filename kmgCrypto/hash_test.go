package kmgCrypto

import (
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestMustMd5File(t *testing.T) {
	kmgFile.MustDelete("testFile")
	kmgFile.MustWriteFile("testFile", []byte("1"))

	md5name := MustMd5File("testFile")
	kmgTest.Equal(md5name, "c4ca4238a0b923820dcc509a6f75849b")
	kmgFile.MustWriteFile("testFile", []byte(md5name))

}
