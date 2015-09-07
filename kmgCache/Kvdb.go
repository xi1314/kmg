package kmgCache

import (
	"encoding/json"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgCrypto"
	"github.com/bronze1man/kmg/kmgFile"
	"os"
)

// 如果没有数据会返回nil
func MustKvdbGetBytes(s string) (b []byte) {
	key := kmgCrypto.Md5Hex([]byte(s))
	content, err := kmgFile.ReadFile(kmgConfig.DefaultEnv().PathInTmp("kvdb/" + key))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		panic(err)
	}
	return content
}

func MustKvdbSetBytes(s string, b []byte) {
	key := kmgCrypto.Md5Hex([]byte(s))
	kmgFile.MustWriteFileWithMkdir(kmgConfig.DefaultEnv().PathInTmp("kvdb/"+key), b)
}

// 返回是否找到了数据
func MustKvdbGet(s string, obj interface{}) bool {
	b := MustKvdbGetBytes(s)
	if b == nil {
		return false
	}
	err := json.Unmarshal(b, obj)
	if err != nil {
		panic(err)
	}
	return true
}

func MustKvdbSet(s string, obj interface{}) {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	MustKvdbSetBytes(s, b)
}
