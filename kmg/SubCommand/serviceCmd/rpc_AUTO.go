package serviceCmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bronze1man/kmg/kmgCrypto"
	"github.com/bronze1man/kmg/kmgLog"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"net/http"
)

type Client_ServiceRpc struct {
	RemoteUrl string // http://kmg.org:1234/
	Psk       *[32]byte
}

//server
func ListenAndServe_ServiceRpc(addr string, obj *ServiceRpc, psk *[32]byte) (closer func() error) {
	s := NewServer_ServiceRpc(obj, psk)
	return kmgHttp.MustGoHttpAsyncListenAndServeWithCloser(addr, s)
}
func NewServer_ServiceRpc(obj *ServiceRpc, psk *[32]byte) http.Handler {
	return &generateServer_ServiceRpc{
		obj: obj,
		psk: psk,
	}
}
func NewClient_ServiceRpc(RemoteUrl string, Psk *[32]byte) *Client_ServiceRpc {
	return &Client_ServiceRpc{RemoteUrl: RemoteUrl, Psk: Psk}
}

type generateServer_ServiceRpc struct {
	obj *ServiceRpc
	psk *[32]byte
}

// http-json-api v1
// 1.数据传输使用psk加密,明文不泄漏信息
// 2.使用json序列化信息
// 3.只有部分api
func (s *generateServer_ServiceRpc) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	b1, err := kmgHttp.RequestReadAllBody(req)
	if err != nil {
		http.Error(w, "error 1", 400)
		kmgLog.Log("InfoServerError", err.Error(), kmgHttp.NewLogStruct(req))
		return
	}
	if s.psk != nil {
		//解密
		b1, err = kmgCrypto.CompressAndEncryptBytesDecodeV2(s.psk, b1)
		if err != nil {
			http.Error(w, "error 2", 400)
			kmgLog.Log("InfoServerError", err.Error(), kmgHttp.NewLogStruct(req))
			return
		}
	}
	outBuf, err := s.handleApiV1(b1)
	if err != nil {
		kmgLog.Log("InfoServerError", err.Error(), kmgHttp.NewLogStruct(req))
		outBuf = append([]byte{1}, err.Error()...) // error
	} else {
		outBuf = append([]byte{2}, outBuf...) // success
	}
	if s.psk != nil {
		//加密
		outBuf = kmgCrypto.CompressAndEncryptBytesEncodeV2(s.psk, outBuf)
	}
	w.WriteHeader(200)
	w.Header().Set("Content-type", "image/jpeg")
	w.Write(outBuf)
}
func (c *Client_ServiceRpc) sendRequest(apiName string, inData interface{}, outData interface{}) (err error) {
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
	if c.Psk != nil {
		inByte = kmgCrypto.CompressAndEncryptBytesEncodeV2(c.Psk, inByte)
	}
	resp, err := http.Post(c.RemoteUrl, "image/jpeg", bytes.NewBuffer(inByte))
	if err != nil {
		return
	}
	outByte, err := kmgHttp.ResponseReadAllBody(resp)
	if err != nil {
		return
	}
	if c.Psk != nil {
		outByte, err = kmgCrypto.CompressAndEncryptBytesDecodeV2(c.Psk, outByte)
		if err != nil {
			return
		}
	}
	if len(outByte) == 0 {
		return errors.New("len(outByte)==0")
	}
	switch outByte[0] {
	case 1: //error
		return errors.New(string(outByte[1:]))
	case 2: //success
		return json.Unmarshal(outByte[1:], outData)
	default:
		return fmt.Errorf("httpjsonApi protocol error 1 %d", outByte[0])
	}
}
func (s *generateServer_ServiceRpc) handleApiV1(inBuf []byte) (outBuf []byte, err error) {
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
	case "Send":
		var Err error
		reqData := &struct {
			Status StartStatus
		}{}
		Err = json.Unmarshal(b2, reqData)
		if Err != nil {
			return nil, Err
		}
		s.obj.Send(reqData.Status)
		return json.Marshal(struct {
		}{})
	}
	return nil, fmt.Errorf("api %s not found", name)
}
func (c *Client_ServiceRpc) Send(Status StartStatus) (Err error) {
	reqData := &struct {
		Status StartStatus
	}{
		Status: Status,
	}
	respData := &struct {
	}{}
	Err = c.sendRequest("Send", reqData, &respData)
	return Err
}

var gClient_ServiceRpc *Client_ServiceRpc

// 全局函数,请先设置客户端的地址,再获取全局客户端,此处不能并发调用
func SetClient_ServiceRpcConfig(RemoteAddr string, Psk *[32]byte) {
	gClient_ServiceRpc = NewClient_ServiceRpc(RemoteAddr, Psk)
}
func GetClient_ServiceRpc() *Client_ServiceRpc {
	return gClient_ServiceRpc
}
