package kmgCrypto

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"math/big"
)

//使用和rsa.SignPKCS1v15相同的算法,但是返回解密后的数据
func RsaPublicDecryptPKCS1v15(pub *rsa.PublicKey, enc []byte) (data []byte, err error) {
	c := new(big.Int).SetBytes(enc)
	m := RsaEncrypt(new(big.Int), pub, c)
	data = m.Bytes()
	//还原数据长度
	for i := 2; i < len(data); i++ {
		if data[i] == 0 {
			return data[i+1:], nil
		}
	}
	return nil, ErrDecryption
	//e := big.NewInt(int64(pub.E))
	//c.Exp(m, e, pub.N)
}

func RsaEncrypt(c *big.Int, pub *rsa.PublicKey, m *big.Int) *big.Int {
	e := big.NewInt(int64(pub.E))
	c.Exp(m, e, pub.N)
	return c
}

var (
	ErrInputSize           = errors.New("input size too large")
	ErrEncryption          = errors.New("encryption error")
	ErrDecryption          = errors.New("decryption error")
	ErrNotRsaTypePublicKey = errors.New("public key is not rsa type.")
)

func RsaPrivateEncryptPKCS1v15(priv *rsa.PrivateKey, data []byte) (enc []byte, err error) {
	return rsa.SignPKCS1v15(nil, priv, crypto.Hash(0), data)
	/*
		k := (priv.N.BitLen() + 7) / 8
		tLen := len(data)
		// rfc2313, section 8:
		// The length of the data D shall not be more than k-11 octets
		if tLen > k-11 {
			err = ErrInputSize
			return
		}
		em := make([]byte, k)
		em[1] = 1
		for i := 2; i < k-tLen-1; i++ {
			em[i] = 0xff
		}
		copy(em[k-tLen:k], data)
		c := new(big.Int).SetBytes(em)
		if c.Cmp(priv.N) > 0 {
			err = ErrEncryption
			return
		}
		var m *big.Int
		var ir *big.Int
		if priv.Precomputed.Dp == nil {
			m = new(big.Int).Exp(c, priv.D, priv.N)
		} else {
			// We have the precalculated values needed for the CRT.
			m = new(big.Int).Exp(c, priv.Precomputed.Dp, priv.Primes[0])
			m2 := new(big.Int).Exp(c, priv.Precomputed.Dq, priv.Primes[1])
			m.Sub(m, m2)
			if m.Sign() < 0 {
				m.Add(m, priv.Primes[0])
			}
			m.Mul(m, priv.Precomputed.Qinv)
			m.Mod(m, priv.Primes[0])
			m.Mul(m, priv.Primes[1])
			m.Add(m, m2)

			for i, values := range priv.Precomputed.CRTValues {
				prime := priv.Primes[2+i]
				m2.Exp(c, values.Exp, prime)
				m2.Sub(m2, m)
				m2.Mul(m2, values.Coeff)
				m2.Mod(m2, prime)
				if m2.Sign() < 0 {
					m2.Add(m2, prime)
				}
				m2.Mul(m2, values.R)
				m.Add(m, m2)
			}
		}

		if ir != nil {
			// Unblind.
			m.Mul(m, ir)
			m.Mod(m, priv.N)
		}
		enc = m.Bytes()
		return
	*/
}

func RsaParseOpensslPrivateKey(b []byte) (key *rsa.PrivateKey, err error) {
	block, _ := pem.Decode(b)
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
func RsaParseOpensslPublicKey(b []byte) (pub *rsa.PublicKey, err error) {
	block, _ := pem.Decode(b)
	pk, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	pub, ok := pk.(*rsa.PublicKey)
	if !ok {
		return nil, ErrNotRsaTypePublicKey
	}
	return pub, nil
}

func RsaOpensslSign(pub *rsa.PrivateKey, h crypto.Hash, msg []byte) (s []byte, err error) {
	h1 := h.New()
	h1.Write(msg)
	digest := h1.Sum(nil)
	s, err = rsa.SignPKCS1v15(nil, pub, h, digest)
	return
}

//这个接口应该和php版的openssl_verify在使用rsa公钥的时候有完全相同的输入输出,加密的坑简直太多了..
//msg是需要验证签名的消息,sig是签名之后生成的
func RsaOpensslVerify(pub *rsa.PublicKey, h crypto.Hash, msg []byte, sig []byte) (err error) {
	h1 := h.New()
	h1.Write(msg)
	digest := h1.Sum(nil)
	err = rsa.VerifyPKCS1v15(pub, h, digest, sig)
	return
}

//支付宝签名的默认方式,
// 放在这里主要用于文档该功能如何实现,
// 并且提供一种签名的方式
// 读入 PKCS1格式私钥字符串 , 需要签名的数据
// 返回 签名后的数据
func RsaWithSha1PKCS1OpensslSignBase64(PKCS1privateKey []byte, msg []byte) (s string, err error) {
	block, _ := pem.Decode(PKCS1privateKey)
	private, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}
	b, err := RsaOpensslSign(private, crypto.SHA1, msg)
	if err != nil {
		return
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

//有时候会得到一个没有头没有尾的publicKey,使用这个函数加上 -----BEGIN PUBLIC KEY----- -----END PUBLIC KEY-----
func PemAddStartEnd(in string) []byte {
	return []byte("-----BEGIN PUBLIC KEY-----\n" + in + "\n-----END PUBLIC KEY-----")
}
