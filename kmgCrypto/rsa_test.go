package kmgCrypto

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/bronze1man/kmg/kmgTest"
	//"encoding/pem"
	//"fmt"
	//"crypto/x509"
	"crypto"

	"github.com/bronze1man/kmg/encoding/kmgBase64"
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

func TestRsaOpensslVerify(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	var err error
	privateKey := []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDA4E8H2qksOnCSoBkq+HH3Dcu0/iWt3iNcpC/BCg0F8tnMhF1Q
OQ98cRUM8eeI9h+S6g/5UmO4hBKMOP3vg/u7kI0ujrCN1RXpsrTbWaqry/xTDgTM
8HkKkNhRSyDJWJVye0mPgbvVnx76en+K6LLzDaQH8yKI/dbswSq65XFcIwIDAQAB
AoGAU+uFF3LBdtf6kSGNsc+lrovXHWoTNOJZWn6ptIFOB0+SClVxUG1zWn7NXPOH
/WSxejfTOXTqpKb6dv55JpSzmzf8fZphVE9Dfr8pU68x8z5ft4yv314qLXFDkNgl
MeQht4n6mo1426dyoOcCfmWc5r7LQCi7WmOsKvATe3nzk/kCQQDp1gyDIVAbUvwe
tpsxZpAd3jLD49OVHUIy2eYGzZZLK3rA1uNWWZGsjrJQvfGf+mW+/zeUMYPBpk0B
XYqlgHJNAkEA0yhhu/2SPJYxIS9umCry1mj7rwji5O2qVSssChFyOctcbysbNJLH
qoF7wumr9PAjjWFWdmCzzEJyxMMurL3gLwJBAIEoeNrJQL0G9jlktY3wz7Ofsrye
j5Syh4kc8EBbuCMnDfOL/iAI8zyzyOxuLhMmNKLtx140h0kkOS6C430M2JUCQCnM
a5RX/JOrs2v7RKwwjENvIqsiWi+w8C/NzPjtPSw9mj2TTd5ZU9bnrMUHlnd09cSt
yPzD5bOAT9GtRVcCexcCQBxXHRleikPTJC90GqrC2l6erYJaiSOtv6QYIh0SEDVm
1o6Whw4FEHUPqMW0Z5PobPFiEQT+fFR02xU3NJrjYy0=
-----END RSA PRIVATE KEY-----`)
	publicKey := []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDA4E8H2qksOnCSoBkq+HH3Dcu0/
iWt3iNcpC/BCg0F8tnMhF1QOQ98cRUM8eeI9h+S6g/5UmO4hBKMOP3vg/u7kI0ujr
CN1RXpsrTbWaqry/xTDgTM8HkKkNhRSyDJWJVye0mPgbvVnx76en+K6LLzDaQH8yK
I/dbswSq65XFcIwIDAQAB
-----END PUBLIC KEY-----`)
	signed := kmgBase64.MustStdBase64DecodeStringToByte(`AqDW/m+aGn2kFo54Bt5XnXniBDtCxmPS6FMfHrLizh7d4jgnz4LbwBfRvXywI6HEKgr7Vk37duTM8P+XqmT+uQU2R1h4nRwOf2fCstXmgeD3qGk/XI+XMafgMkTnV/B9dOXpdUbxEpL1fDhmo7A6J0rcJotG7TP7i1zcvY4oiXk=`)
	msg := []byte("this is a test!")
	rsaKey, err := RsaParseOpensslPrivateKey(privateKey)
	s, err := RsaOpensslSign(rsaKey, crypto.SHA1, msg)
	t.Equal(err, nil)
	t.Equal(s, signed)

	rsapk, err := RsaParseOpensslPublicKey(publicKey)
	t.Equal(err, nil)
	err = RsaOpensslVerify(rsapk, crypto.SHA1, msg, signed)
	t.Equal(err, nil)
}
