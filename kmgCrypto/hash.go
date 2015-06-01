package kmgCrypto

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"github.com/bronze1man/kmg/kmgFile"
)

func Sha256Hex(data []byte) string {
	out := sha256.Sum256(data)
	return hex.EncodeToString(out[:])
}

func Sha512Hex(data []byte) string {
	out := sha512.Sum512(data)
	return hex.EncodeToString(out[:])
}
func Sha512HexFromString(data string) string {
	out := sha512.Sum512([]byte(data))
	return hex.EncodeToString(out[:])
}

//小写hex
func Md5Hex(data []byte) string {
	hash := md5.New()
	hash.Write(data)
	return hex.EncodeToString(hash.Sum(nil))
}

//小写hex
func Md5HexFromString(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

//获得文件的MD5值
func MustMd5File(path string) string {
	content := kmgFile.MustReadFile(path)
	return Md5Hex(content)
}
