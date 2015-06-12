package kmgCompress

import (
	"bytes"
	"compress/flate"
	"io/ioutil"
)

//flate压缩,panic+buffer实现
// 会由于
func FlateMustCompress(inb []byte) (outb []byte) {
	buf := &bytes.Buffer{}
	w, err := flate.NewWriter(buf, -1)
	if err != nil {
		panic(err) //只会由于加密level设置错误而报错 但是上面明显是对的.
	}
	_, err = w.Write(inb)
	if err != nil {
		w.Close()
		panic(err)
	}
	err = w.Close()
	if err != nil {
		panic(err) // 全部代码看过,只会由于w.Write错误而报错, 但是明显w.Write 不会报错.
	}
	return buf.Bytes()
}

func FlateMustUnCompress(inb []byte) (outb []byte) {
	buf := bytes.NewBuffer(inb)
	reader := flate.NewReader(buf)
	outb, err := ioutil.ReadAll(reader)
	if err != nil {
		reader.Close()
		panic(err)
	}
	err = reader.Close()
	if err != nil {
		panic(err)
	}
	return outb
}

func FlateUnCompress(inb []byte) (outb []byte, err error) {
	buf := bytes.NewBuffer(inb)
	reader := flate.NewReader(buf)
	outb, err = ioutil.ReadAll(reader)
	if err != nil {
		reader.Close()
		return nil, err
	}
	err = reader.Close()
	if err != nil {
		return
	}
	return outb, nil
}
