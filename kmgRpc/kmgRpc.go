package kmgRpc

import (
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgTextTemplate"
	"go/format"
	"strings"
	//"fmt"
)

type GenerateRequest struct {
	Object               interface{} //需要分析的对象
	ObjectName           string      //TODO 分析出对象名称来.
	OutFilePath          string      //生成的文件路径
	OutPackageImportPath string      //生成的package的importPath
	Key                  [32]byte    //密钥
}

//生成代码
// 会把Object上面所有的公开函数都拿去生成一遍
func MustGenerateCode(req GenerateRequest) {
	config := reflectToTplConfig(req)
	outB := tplGenerateCode(config)
	//fmt.Println(string(outB))
	outB1, err := format.Source(outB)
	if err == nil {
		outB = outB1
	}
	kmgFile.MustWriteFileWithMkdir(req.OutFilePath, outB)
	return
}

type tplConfig struct {
	OutPackageName string          //生成的package的名字 testPackage
	OutKeyBase64   string          //生成的key的base64的值
	ObjectName     string          //对象名字	如 Demo
	ObjectTypeStr  string          //对象的类型表示	如 *Demo
	ImportPathMap  map[string]bool //ImportPath列表
	ApiList        []Api           //api列表
}

func (conf *tplConfig) mergeImportPath(importPathList []string) {
	for _, importPath := range importPathList {
		conf.ImportPathMap[importPath] = true
	}
}

type Api struct {
	Name        string                 //在这个系统里面的名字
	InArgsList  []ArgumentNameTypePair //输入变量列表
	OutArgsList []ArgumentNameTypePair //输出变量列表
}

func (api Api) GetOutArgsListWithoutError() []ArgumentNameTypePair {
	out := make([]ArgumentNameTypePair, 0, len(api.OutArgsList))
	for _, pair := range api.OutArgsList {
		if pair.ObjectTypeStr == "error" {
			continue
		}
		out = append(out, pair)
	}
	return out
}
func (api Api) GetOutArgsNameListForAssign() string {
	nameList := []string{}
	for _, pair := range api.OutArgsList {
		nameList = append(nameList, pair.Name)
	}
	return strings.Join(nameList, ",")
}

func (api Api) HasReturnArgument() bool {
	return len(api.OutArgsList) > 0
}

func (api Api) GetClientOutArgument() []ArgumentNameTypePair {
	for _, pair := range api.OutArgsList {
		if pair.ObjectTypeStr == "error" {
			return api.OutArgsList
		}
	}
	return append(api.OutArgsList, ArgumentNameTypePair{
		Name:          "err",
		ObjectTypeStr: "error",
	})
}

// TODO 下一个版本不要这个hook了,复杂度太高
func (api Api) IsOutExpendToOneArgument() bool {
	return len(api.OutArgsList) == 2 &&
		api.OutArgsList[0].Name == "Response" &&
		api.OutArgsList[1].ObjectTypeStr == "error"
}

type ArgumentNameTypePair struct {
	Name          string
	ObjectTypeStr string
}

/*
func tplGenerateCode(config tplConfig) []byte {
	return kmgTextTemplate.MustRenderToByte(`package {{.OutPackageName}}

import (
	{{range $index, $element := .ImportPathMap}}"{{$index}}"
	{{end}}
)

//server
func ListenAndServe_{{.ObjectName}}(addr string, obj {{.ObjectTypeStr}}) {
	s := &generateServer_{{.ObjectName}}{
		obj: obj,
	}
	err := http.ListenAndServe(addr, s)
	if err != nil {
		panic(err)
	}
}

func NewServer_{{.ObjectName}}(obj {{.ObjectTypeStr}}) http.Handler {
	return &generateServer_{{.ObjectName}}{
		obj: obj,
	}
}

func NewClient_{{.ObjectName}}(RemoteUrl string) *Client_{{.ObjectName}} {
	return &Client_{{.ObjectName}}{RemoteUrl: RemoteUrl}
}

//client
// 信息服务器的客户端.
// httpjson api v1 client used for monitor to check that the server is good.
type Client_{{.ObjectName}} struct {
	RemoteUrl string //只有主机和地址
}


var kmgRpc_{{.ObjectName}}_encryptKey = kmgBase64.MustStdBase64DecodeString("{{.OutKeyBase64}}")

const (
	kmgRpc_{{.ObjectName}}_ResponseCodeSuccess byte = 1
	kmgRpc_{{.ObjectName}}_ResponseCodeError   byte = 2
)



type generateServer_{{.ObjectName}} struct {
	obj {{.ObjectTypeStr}}
}

// http-json-api v1
// 1.数据传输使用psk加密,明文不泄漏信息
// 2.使用json序列化信息
// 3.只有部分api
func (s *generateServer_{{.ObjectName}}) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	b1, err := kmgHttp.RequestReadAllBody(req)
	if err != nil {
		http.Error(w, "error 1", 400)
		kmgLog.Log("InfoServerError", err.Error(), kmgHttp.NewLogStruct(req))
		return
	}

	//解密
	b1, err = kmgCrypto.AesCbcPKCS7PaddingDecrypt(b1, kmgRpc_{{.ObjectName}}_encryptKey)
	if err != nil {
		http.Error(w, "error 2", 400)
		kmgLog.Log("InfoServerError", err.Error(), kmgHttp.NewLogStruct(req))
		return
	}
	outBuf, err := s.handleApiV1(b1)
	if err != nil {
		kmgLog.Log("InfoServerError", err.Error(), kmgHttp.NewLogStruct(req))
		outBuf = append([]byte{kmgRpc_{{.ObjectName}}_ResponseCodeError}, err.Error()...)
	} else {
		outBuf = append([]byte{kmgRpc_{{.ObjectName}}_ResponseCodeSuccess}, outBuf...)
	}
	//加密
	outBuf = kmgCrypto.AesCbcPKCS7PaddingEncrypt(outBuf, kmgRpc_{{.ObjectName}}_encryptKey)
	w.WriteHeader(200)
	w.Header().Set("Content-type", "image/jpeg")
	w.Write(outBuf)
}

func (c *Client_{{.ObjectName}}) sendRequest(apiName string, inData interface{}, outData interface{}) (err error) {
	inDataByte, err := json.Marshal(inData)
	if err != nil {
		return
	}
	if len(apiName) > 255 {
		return errors.New("len(apiName)>255")
	}
	inByte := []byte{byte(len(apiName))}
	inByte = append(inByte, []byte(apiName)...)
	inByte = append(inByte, inDataByte...)
	inByte = kmgCrypto.AesCbcPKCS7PaddingEncrypt(inByte, kmgRpc_{{.ObjectName}}_encryptKey)

	resp, err := http.Post(c.RemoteUrl, "image/jpeg", bytes.NewBuffer(inByte))
	if err != nil {
		return
	}
	outByte, err := kmgHttp.ResponseReadAllBody(resp)
	if err != nil {
		return
	}
	outByte, err = kmgCrypto.AesCbcPKCS7PaddingDecrypt(outByte, kmgRpc_{{.ObjectName}}_encryptKey)
	if err != nil {
		return
	}
	if len(outByte) == 0 {
		return errors.New("len(outByte)==0")
	}
	switch outByte[0] {
	case kmgRpc_{{.ObjectName}}_ResponseCodeError:
		return errors.New(string(outByte[1:]))
	case kmgRpc_{{.ObjectName}}_ResponseCodeSuccess:
		return json.Unmarshal(outByte[1:], outData)
	default:
		return fmt.Errorf("httpjsonApi protocol error 1 %d", outByte[0])
	}
}



func (s *generateServer_{{.ObjectName}}) handleApiV1(inBuf []byte) (outBuf []byte, err error) {
	//从此处开始协议正确了,换一种返回方式
	// 1 byte api name len apiNameLen
	// apiNameLen byte api name
	// xx byte json encode of request as struct.
	if len(inBuf) < 2 {
		return nil, fmt.Errorf("len(b1)<2")
	}
	nameLength := inBuf[0]
	if len(inBuf) < int(nameLength)+1 {
		return nil, fmt.Errorf("len(b1)<nameLength+1")
	}
	name := string(inBuf[1 : int(nameLength)+1])
	b2 := inBuf[nameLength+1:]

	switch name {
	{{range .ApiList}}
	case "{{.Name}}":
		{{range .GetOutArgsListWithoutError}}var {{.Name}} {{.ObjectTypeStr}}
		{{end}}
		var err error
		reqData := &struct {
			{{range .InArgsList}}{{.Name}} {{.ObjectTypeStr}}
			{{end}}
		}{}
		err = json.Unmarshal(b2, reqData)
		if err != nil {
			return nil, err
		}
		{{if .HasReturnArgument}}
			{{.GetOutArgsNameListForAssign}} = s.obj.{{.Name}}( {{range .InArgsList}} reqData.{{.Name}},{{end}} )
			if err != nil {
				return nil, err
			}
		{{else}}
			s.obj.{{.Name}}( {{range .InArgsList}} reqData.{{.Name}},{{end}} )
		{{end}}
		{{if .IsOutExpendToOneArgument}}
			return json.Marshal(Response)
		{{else}}
			return json.Marshal(struct {
				{{range .GetOutArgsListWithoutError}}{{.Name}} {{.ObjectTypeStr}}
				{{end}}
			}{
				{{range .GetOutArgsListWithoutError}}{{.Name}}:{{.Name}},
				{{end}}
			})
		{{end}}

	{{end}}
	}
	return nil, fmt.Errorf("api %s not found", name)
}

{{range .ApiList}}
func (c *Client_{{$.ObjectName}}) {{.Name}}( {{range .InArgsList}} {{.Name}} {{.ObjectTypeStr}}, {{end}} ) ({{range .GetClientOutArgument}} {{.Name}} {{.ObjectTypeStr}}, {{end}}) {
	reqData := &struct {
		{{range .InArgsList}}{{.Name}} {{.ObjectTypeStr}}
		{{end}}
	}{
		{{range .InArgsList}}{{.Name}}: {{.Name}},
		{{end}}
	}
	{{if .IsOutExpendToOneArgument}}
		var respData {{(index .OutArgsList 0).ObjectTypeStr}}
		err = c.sendRequest("{{.Name}}", reqData, &respData)
		return respData,err
	{{else}}
		respData := &struct {
			{{range .GetOutArgsListWithoutError}}{{.Name}} {{.ObjectTypeStr}}
			{{end}}
		}{}
		err = c.sendRequest("{{.Name}}", reqData, &respData)
		return {{range .GetOutArgsListWithoutError}}respData.{{.Name}},{{end}} err
	{{end}}
}
{{end}}


`, config)
}
*/
