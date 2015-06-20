package kmgJson

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/bronze1man/kmg/typeTransform"
	"strings"
)

func ReadFile(path string, obj interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, obj)
}

func MustReadFile(path string, obj interface{}) {
	err := ReadFile(path, obj)
	if err != nil {
		panic(err)
	}
}

func MustWriteFileIndent(path string, obj interface{}) {
	output, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(path, output, os.FileMode(0777))
	if err != nil {
		panic(err)
	}
}

//读取json文件,并修正json的类型问题(map key 必须是string的问题)
func ReadFileTypeFix(path string, obj interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var iobj interface{}
	err = json.Unmarshal(b, &iobj)
	if err != nil {
		return err
	}
	return typeTransform.Transform(iobj, obj)
}

func WriteFile(path string, obj interface{}) (err error) {
	out, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, out, os.FileMode(0777))
}

//写入json文件,并修正json的类型问题(map key 必须是string的问题)
func WriteFileTypeFix(path string, obj interface{}) (err error) {
	//a simple work around
	obj, err = TypeFixWhenMarshal(obj)
	if err != nil {
		return
	}
	outByte, err := json.Marshal(obj)
	if err != nil {
		return
	}
	return ioutil.WriteFile(path, outByte, os.FileMode(0777))
}

func UnmarshalNoType(r []byte) (interface{}, error) {
	var obj interface{}
	err := json.Unmarshal(r, &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func MustUnmarshal(r []byte, obj interface{}) {
	err := json.Unmarshal(r, &obj)
	if err != nil {
		panic(err)
	}
	return
}

func MustUnmarshalIgnoreEmptyString(jsonStr string, obj interface{}) {
	if jsonStr == "" {
		return
	}
	err := json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		panic(err)
	}
	return
}

func MustUnmarshalToMap(r []byte) (obj map[string]interface{}) {
	err := json.Unmarshal(r, &obj)
	if err != nil {
		panic(err)
	}
	return obj
}

// for debug to inspect content in obj
func MustMarshalIndentToString(obj interface{}) string {
	output, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(output)
}

var htmlUnescapeReplacer = strings.NewReplacer(`\u003c`, "<", `\u003e`, ">", `\u0026`, "&")

func MarshalIndent(obj interface{}) ([]byte, error) {
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return nil, err
	}
	return []byte(htmlUnescapeReplacer.Replace(string(b))), nil
}

func MustMarshal(obj interface{}) []byte {
	output, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return output
}

func MustMarshalToString(obj interface{}) string {
	output, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(output)
}
