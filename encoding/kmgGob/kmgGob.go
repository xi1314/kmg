package kmgGob

import (
	"bytes"
	"encoding/gob"
	"io"
	"os"

	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgIo"
)

func WriteFile(path string, obj interface{}) (err error) {
	err = kmgFile.MkdirForFile(path)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0666))
	if err != nil {
		return
	}
	defer f.Close()
	encoder := gob.NewEncoder(f)
	return encoder.Encode(obj)
}

func MustWriteFile(path string, obj interface{}) {
	err := WriteFile(path, obj)
	if err != nil {
		panic(err)
	}
}

func ReadFile(path string, obj interface{}) (err error) {
	f, err := os.OpenFile(path, os.O_RDONLY, os.FileMode(0666))
	if err != nil {
		return
	}
	defer f.Close()
	encoder := gob.NewDecoder(f)
	return encoder.Decode(obj)
}

func MustReadFile(path string, obj interface{}) {
	err := ReadFile(path, obj)
	if err != nil {
		panic(err)
	}
}

func Marshal(obj interface{}) (out []byte, err error) {
	b := &bytes.Buffer{}
	encoder := gob.NewEncoder(b)
	err = encoder.Encode(obj)
	if err != nil {
		return
	}
	return b.Bytes(), nil
}

func MustMarshal(obj interface{}) (out []byte) {
	out, err := Marshal(obj)
	if err != nil {
		panic(err)
	}
	return out
}
func Unmarshal(data []byte, obj interface{}) (err error) {
	b := bytes.NewBuffer(data)
	encoder := gob.NewDecoder(b)
	err = encoder.Decode(obj)
	return
}
func MustUnmarshal(data []byte, obj interface{}) {
	err := Unmarshal(data, obj)
	if err != nil {
		panic(err)
	}
	return
}

// 创建一个没有中间buf的gobDecode ,
// 系统已有的gob.NewDecoder, 会在Reader不是一个ByteReader时, 使用bufio 进行读取,
// 根据 gob.NewDecoder 的代码显示 使用singleReader可以解决问题
func NewNonBufDecode(r io.Reader) *gob.Decoder {
	return gob.NewDecoder(kmgIo.NewSingleByteReader(r))
}
