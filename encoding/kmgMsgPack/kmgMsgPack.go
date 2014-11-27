package kmgMsgPack

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/ugorji/go/codec"
)

func WriteFile(path string, obj interface{}) (err error) {
	mh := codec.MsgpackHandle{}
	mh.AsSymbols = codec.AsSymbolNone
	mh.RawToString = true
	buf := &bytes.Buffer{}
	encoder := codec.NewEncoder(buf, &mh)
	err = encoder.Encode(obj)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path, buf.Bytes(), os.FileMode(0777))
	return
}
