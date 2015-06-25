package kmgRpc

import (
	"bytes"
)

func tplGenerateCode(config tplConfig) []byte {
	var _buf bytes.Buffer
	_buf.WriteString(`
package `)
	_buf.WriteString(config.OutPackageName)
	_buf.WriteString(`

import (
    `)
	for _, importPkg := range config.ImportPathMap {
		_buf.WriteString(`
        "`)
		_buf.WriteString(importPkg)
		_buf.WriteString(`"
    `)
	}
	_buf.WriteString(`
)

//server
func ListenAndServe_`)
	_buf.WriteString(config.ObjectName)
	_buf.WriteString(`(addr string, obj `)
	_buf.WriteString(config.ObjectTypeStr)
	_buf.WriteString(`) {
	s := &generateServer_`)
	_buf.WriteString(config.ObjectName)
	_buf.WriteString(`{
		obj: obj,
	}
	err := http.ListenAndServe(addr, s)
	if err != nil {
		panic(err)
	}
}

func NewServer_`)
	_buf.WriteString(config.ObjectName)
	_buf.WriteString(`(obj `)
	_buf.WriteString(config.ObjectTypeStr)
	_buf.WriteString(`) http.Handler {
	return &generateServer_`)
	_buf.WriteString(config.ObjectName)
	_buf.WriteString(`{
		obj: obj,
	}
}

func NewClient_`)
	_buf.WriteString(config.ObjectName)
	_buf.WriteString(`(RemoteUrl string) *Client_`)
	_buf.WriteString(config.ObjectName)
	_buf.WriteString(` {
	return &Client_`)
	_buf.WriteString(config.ObjectName)
	_buf.WriteString(`{RemoteUrl: RemoteUrl}
}

//client
// 信息服务器的客户端.
// httpjson api v1 client used for monitor to check that the server is good.
type Client_`)
	_buf.WriteString(config.ObjectName)
	_buf.WriteString(` struct {
	RemoteUrl string //只有主机和地址
}


var kmgRpc_`)
	_buf.WriteString(config.ObjectName)
	_buf.WriteString(`_encryptKey = kmgBase64.MustStdBase64DecodeString("`)
	_buf.WriteString(config.OutKeyBase64)
	_buf.WriteString(`")

const (
	kmgRpc_`)
	_buf.WriteString(config.ObjectName)
	_buf.WriteString(`_ResponseCodeSuccess byte = 1
	kmgRpc_`)
	_buf.WriteString(config.ObjectName)
	_buf.WriteString(`_ResponseCodeError   byte = 2
)



type generateServer_`)
	_buf.WriteString(config.ObjectName)
	_buf.WriteString(` struct {
	obj `)
	_buf.WriteString(config.ObjectTypeStr)
	_buf.WriteString(`
}

// http-json-api v1
// 1.数据传输使用psk加密,明文不泄漏信息
// 2.使用json序列化信息
// 3.只有部分api
func (s *generateServer_`)
	_buf.WriteString(config.ObjectName)
	_buf.WriteString(`) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	b1, err := kmgHttp.RequestReadAllBody(req)
	if err != nil {
		http.Error(w, "error 1", 400)
		kmgLog.Log("InfoServerError", err.Error(), kmgHttp.NewLogStruct(req))
		return
	}

	//解密
	b1, err = kmgCrypto.AesCbcPKCS7PaddingDecrypt(b1, kmgRpc_`)
	_buf.WriteString(config.ObjectName)
	_buf.WriteString(`_encryptKey)
	if err != nil {
		http.Error(w, "error 2", 400)
		kmgLog.Log("InfoServerError", err.Error(), kmgHttp.NewLogStruct(req))
		return
	}
	outBuf, err := s.handleApiV1(b1)
	if err != nil {
		kmgLog.Log("InfoServerError", err.Error(), kmgHttp.NewLogStruct(req))
		outBuf = append([]byte{kmgRpc_`)
	_buf.WriteString(config.ObjectName)
	_buf.WriteString(`_ResponseCodeError}, err.Error()...)
	} else {
		outBuf = append([]byte{kmgRpc_`)
	_buf.WriteString(config.ObjectName)
	_buf.WriteString(`_ResponseCodeSuccess}, outBuf...)
	}
	//加密
	outBuf = kmgCrypto.AesCbcPKCS7PaddingEncrypt(outBuf, kmgRpc_`)
	_buf.WriteString(config.ObjectName)
	_buf.WriteString(`_encryptKey)
	w.WriteHeader(200)
	w.Header().Set("Content-type", "image/jpeg")
	w.Write(outBuf)
}

func (c *Client_`)
	_buf.WriteString(config.ObjectName)
	_buf.WriteString(`) sendRequest(apiName string, inData interface{}, outData interface{}) (err error) {
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
	inByte = kmgCrypto.AesCbcPKCS7PaddingEncrypt(inByte, kmgRpc_`)
	_buf.WriteString(config.ObjectName)
	_buf.WriteString(`_encryptKey)

	resp, err := http.Post(c.RemoteUrl, "image/jpeg", bytes.NewBuffer(inByte))
	if err != nil {
		return
	}
	outByte, err := kmgHttp.ResponseReadAllBody(resp)
	if err != nil {
		return
	}
	outByte, err = kmgCrypto.AesCbcPKCS7PaddingDecrypt(outByte, kmgRpc_`)
	_buf.WriteString(config.ObjectName)
	_buf.WriteString(`_encryptKey)
	if err != nil {
		return
	}
	if len(outByte) == 0 {
		return errors.New("len(outByte)==0")
	}
	switch outByte[0] {
	case kmgRpc_`)
	_buf.WriteString(config.ObjectName)
	_buf.WriteString(`_ResponseCodeError:
		return errors.New(string(outByte[1:]))
	case kmgRpc_`)
	_buf.WriteString(config.ObjectName)
	_buf.WriteString(`_ResponseCodeSuccess:
		return json.Unmarshal(outByte[1:], outData)
	default:
		return fmt.Errorf("httpjsonApi protocol error 1 %d", outByte[0])
	}
}



func (s *generateServer_`)
	_buf.WriteString(config.ObjectName)
	_buf.WriteString(`) handleApiV1(inBuf []byte) (outBuf []byte, err error) {
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
	`)
	for _, api := range config.ApiList {
		_buf.WriteString(`
	case "`)
		_buf.WriteString(api.Name)
		_buf.WriteString(`":
	    `)
		for _, args := range api.GetOutArgsListWithoutError() {
			_buf.WriteString(`
	       var `)
			_buf.WriteString(args.Name)
			_buf.WriteString(` `)
			_buf.WriteString(args.ObjectTypeStr)
			_buf.WriteString(`
	    `)
		}
		_buf.WriteString(`
		var err error
		reqData := &struct {
            `)
		for _, args := range api.InArgsList {
			_buf.WriteString(`
               `)
			_buf.WriteString(args.Name)
			_buf.WriteString(` `)
			_buf.WriteString(args.ObjectTypeStr)
			_buf.WriteString(`
            `)
		}
		_buf.WriteString(`
		}{}
		err = json.Unmarshal(b2, reqData)
		if err != nil {
			return nil, err
		}
		`)
		if api.HasReturnArgument() {
			_buf.WriteString(`
		    `)
			_buf.WriteString(api.GetOutArgsNameListForAssign())
			_buf.WriteString(` = s.obj.`)
			_buf.WriteString(api.Name)
			_buf.WriteString(`(`)
			for _, args := range api.InArgsList {
				_buf.WriteString(` reqData.`)
				_buf.WriteString(args.Name)
				_buf.WriteString(`,`)
			}
			_buf.WriteString(` )
            if err != nil {
                return nil, err
            }
		`)
		} else {
			_buf.WriteString(`
		    s.obj.`)
			_buf.WriteString(api.Name)
			_buf.WriteString(`(`)
			for _, args := range api.InArgsList {
				_buf.WriteString(` reqData.`)
				_buf.WriteString(args.Name)
				_buf.WriteString(`,`)
			}
			_buf.WriteString(` )
		`)
		}
		_buf.WriteString(`
		`)
		if api.IsOutExpendToOneArgument() {
			_buf.WriteString(`
			return json.Marshal(Response)
        `)
		} else {
			_buf.WriteString(`
			return json.Marshal(struct {
			    `)
			for _, arg := range api.GetOutArgsListWithoutError() {
				_buf.WriteString(`
			        `)
				_buf.WriteString(arg.Name)
				_buf.WriteString(` `)
				_buf.WriteString(arg.ObjectTypeStr)
				_buf.WriteString(`
			    `)
			}
			_buf.WriteString(`
			}{
                `)
			for _, arg := range api.GetOutArgsListWithoutError() {
				_buf.WriteString(`
                    `)
				_buf.WriteString(arg.Name)
				_buf.WriteString(`:`)
				_buf.WriteString(arg.Name)
				_buf.WriteString(`
                `)
			}
			_buf.WriteString(`
			})
		`)
		}
		_buf.WriteString(`
	`)
	}
	_buf.WriteString(`
	}
	return nil, fmt.Errorf("api %s not found", name)
}

`)
	for _, api := range config.ApiList {
		_buf.WriteString(`
func (c *Client_`)
		_buf.WriteString(config.ObjectName)
		_buf.WriteString(` ) `)
		_buf.WriteString(api.Name)
		_buf.WriteString(`( `)
		for _, arg := range api.InArgsList {
			_buf.WriteString(arg.Name)
			_buf.WriteString(` `)
			_buf.WriteString(arg.ObjectTypeStr)
			_buf.WriteString(`, `)
		}
		_buf.WriteString(`  ) (`)
		for _, arg := range api.GetClientOutArgument() {
			_buf.WriteString(arg.Name)
			_buf.WriteString(` `)
			_buf.WriteString(arg.ObjectTypeStr)
			_buf.WriteString(`, `)
		}
		_buf.WriteString(` ) {
	reqData := &struct {
	    `)
		for _, arg := range api.InArgsList {
			_buf.WriteString(`
	        `)
			_buf.WriteString(arg.Name)
			_buf.WriteString(` `)
			_buf.WriteString(arg.ObjectTypeStr)
			_buf.WriteString(`
	    `)
		}
		_buf.WriteString(`
	}{
        `)
		for _, arg := range api.InArgsList {
			_buf.WriteString(`
            `)
			_buf.WriteString(arg.Name)
			_buf.WriteString(` `)
			_buf.WriteString(arg.Name)
			_buf.WriteString(`
        `)
		}
		_buf.WriteString(`
	}
	`)
		if api.IsOutExpendToOneArgument() {
			_buf.WriteString(`
	    var respData `)
			_buf.WriteString(api.OutArgsList[0].ObjectTypeStr)
			_buf.WriteString(`
        err = c.sendRequest("`)
			_buf.WriteString(api.Name)
			_buf.WriteString(`", reqData, &respData)
        return respData,err
	`)
		} else {
			_buf.WriteString(`
        respData := &struct {
            `)
			for _, arg := range api.GetOutArgsListWithoutError() {
				_buf.WriteString(`
                `)
				_buf.WriteString(arg.Name)
				_buf.WriteString(` `)
				_buf.WriteString(arg.ObjectTypeStr)
				_buf.WriteString(`
            `)
			}
			_buf.WriteString(`
        }{}
        err = c.sendRequest("`)
			_buf.WriteString(api.Name)
			_buf.WriteString(`", reqData, &respData)
        return `)
			for _, arg := range api.GetOutArgsListWithoutError() {
				_buf.WriteString(`respData.`)
				_buf.WriteString(arg.Name)
				_buf.WriteString(`,`)
			}
			_buf.WriteString(` err
    `)
		}
		_buf.WriteString(`
}
`)
	}
	_buf.WriteString(`
`)
	return _buf.String()
}
