package kmgQiniu

import (
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

var testContext *Context

func TestContextWithFile(ot *testing.T) {
	if testContext == nil {
		ot.Skip("need config testContext todo the real upload download test")
		return
	}
	var err error
	err = testContext.RemovePrefix("/kmgTest/")
	kmgTest.Equal(err, nil)

	err = testContext.UploadFromFile("testFile", "/kmgTest/")
	kmgTest.Equal(err, nil)

	kmgFile.MustDeleteFile("testFile/downloaded.txt")
	err = testContext.DownloadToFile("/kmgTest/1.txt", "testFile/downloaded.txt")
	kmgTest.Equal(err, nil)

	kmgTest.Equal(kmgFile.MustReadFileAll("testFile/downloaded.txt"), []byte("abc"))
}

func TestContextWithBytes(ot *testing.T) {
	if testContext == nil {
		ot.Skip("need config testContext todo the real upload download test")
		return
	}
	var err error
	err = testContext.RemovePrefix("/kmgTest/")
	kmgTest.Equal(err, nil)

	err = testContext.UploadFromBytes("/kmgTest/1.txt", []byte("abc"))
	kmgTest.Equal(err, nil)

	b := testContext.MustDownloadToBytes("/kmgTest/1.txt")
	kmgTest.Equal(b, []byte("abc"))
}

func TestContextWithBytesLeadingSlash(ot *testing.T) {
	if testContext == nil {
		ot.Skip("need config testContext todo the real upload download test")
		return
	}
	var err error
	err = testContext.RemovePrefix("kmgTest/")
	kmgTest.Equal(err, nil)

	err = testContext.UploadFromBytes("kmgTest/1.txt", []byte("abc"))
	kmgTest.Equal(err, nil)

	b := testContext.MustDownloadToBytes("kmgTest/1.txt")
	kmgTest.Equal(b, []byte("abc"))
}
