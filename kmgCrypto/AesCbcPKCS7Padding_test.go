package kmgCrypto

import (
	"bytes"
	"testing"

	"github.com/bronze1man/kmg/encoding/kmgBase64"
	"github.com/bronze1man/kmg/kmgTest"
)

func TestAesCbc(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	in := []byte("123")
	key := kmgBase64.MustStdBase64DecodeString("6VRHJsip3mQ2r1qVI9Gbk7h2v0yfQjqQTbs8hFKUMRo=")
	out := AesCbcPKCS7PaddingEncrypt(in, key)
	out2, err := AesCbcPKCS7PaddingDecrypt(out, key)
	t.Equal(err, nil)
	t.Equal(in, out2)

	out2, err = AesCbcPKCS7PaddingDecrypt(kmgBase64.MustStdBase64DecodeString("BemxQXTV13r3D6dToE0o9qRNNhaNM5VlvZQu6CCcLjc="), key)
	t.Equal(err, nil)
	t.Equal(in, out2)

	for _, in := range [][]byte{
		[]byte(""),
		[]byte("1"),
		[]byte("12"),
		[]byte("123"),
		[]byte("1234"),
		[]byte("12345"),
		[]byte("123456"),
		[]byte("1234567"),
		[]byte("12345678"),
		[]byte("123456789"),
		[]byte("1234567890"),
		[]byte("12345678901"),
		[]byte("123456789012"),
		[]byte("1234567890123"),
		[]byte("12345678901234"),
		[]byte("123456789012345"),
		[]byte("1234567890123456"),
		[]byte("12345678901234567"),
		bytes.Repeat([]byte("123456"), 101),
	} {
		key := kmgBase64.MustStdBase64DecodeString("6VRHJsip3mQ2r1qVI9Gbk7h2v0yfQjqQTbs8hFKUMRo=")
		out := AesCbcPKCS7PaddingEncrypt(in, key)
		out2, err := AesCbcPKCS7PaddingDecrypt(out, key)
		t.Equal(err, nil)
		t.Equal(in, out2)
	}
}
