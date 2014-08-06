package kmgCrypto

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
	//"fmt"
)

func TestRsaPublicDecryptPKCS1v15(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	priv, err := rsa.GenerateKey(rand.Reader, 256)
	t.Equal(err, nil)
	for _, datas := range []string{
		"\000",
		"\000\000",
		"123456789012345678901",
	} {
		data := []byte(datas)
		enc, err := RsaPrivateEncryptPKCS1v15(priv, data)
		t.Equal(err, nil)
		dout, err := RsaPublicDecryptPKCS1v15(&priv.PublicKey, enc)
		t.Equal(err, nil)
		t.Equal(dout, data)
	}

}
