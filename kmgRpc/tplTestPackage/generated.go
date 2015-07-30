package tplTestPackage

import (
	"github.com/bronze1man/kmg/kmgLog"

	"github.com/bronze1man/kmg/kmgNet/kmgHttp"

	"net/http"

	"bytes"

	"encoding/json"

	"errors"

	"fmt"

	"github.com/bronze1man/kmg/kmgCrypto"
)

type Client_Demo struct {
	RemoteUrl string // http://kmg.org:1234/
}

var kmgRpc_Demo_encryptKey = &[32]byte{1, 2}

const (
	kmgRpc_Demo_ResponseCodeSuccess byte = 1
	kmgRpc_Demo_ResponseCodeError   byte = 2
)

//server
func ListenAndServe_Demo(addr string, obj *Demo) (closer func() error) {
	s := &generateServer_Demo{
		obj: obj,
	}
	return kmgHttp.MustGoHttpAsyncListenAndServeWithCloser(addr, s)
}

func NewServer_Demo(obj *Demo) http.Handler {
	return &generateServer_Demo{
		obj: obj,
	}
}

func NewClient_Demo(RemoteUrl string) *Client_Demo {
	return &Client_Demo{RemoteUrl: RemoteUrl}
}

type generateServer_Demo struct {
	obj *Demo
}

// http-json-api v1
// 1.数据传输使用psk加密,明文不泄漏信息
// 2.使用json序列化信息
// 3.只有部分api
func (s *generateServer_Demo) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	b1, err := kmgHttp.RequestReadAllBody(req)
	if err != nil {
		http.Error(w, "error 1", 400)
		kmgLog.Log("InfoServerError", err.Error(), kmgHttp.NewLogStruct(req))
		return
	}

	//解密
	b1, err = kmgCrypto.CompressAndEncryptBytesDecode(kmgRpc_Demo_encryptKey, b1)
	if err != nil {
		http.Error(w, "error 2", 400)
		kmgLog.Log("InfoServerError", err.Error(), kmgHttp.NewLogStruct(req))
		return
	}
	outBuf, err := s.handleApiV1(b1)
	if err != nil {
		kmgLog.Log("InfoServerError", err.Error(), kmgHttp.NewLogStruct(req))
		outBuf = append([]byte{kmgRpc_Demo_ResponseCodeError}, err.Error()...)
	} else {
		outBuf = append([]byte{kmgRpc_Demo_ResponseCodeSuccess}, outBuf...)
	}
	//加密
	outBuf = kmgCrypto.CompressAndEncryptBytesEncode(kmgRpc_Demo_encryptKey, outBuf)
	w.WriteHeader(200)
	w.Header().Set("Content-type", "image/jpeg")
	w.Write(outBuf)
}

func (c *Client_Demo) sendRequest(apiName string, inData interface{}, outData interface{}) (err error) {
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
	inByte = kmgCrypto.CompressAndEncryptBytesEncode(kmgRpc_Demo_encryptKey, inByte)

	resp, err := http.Post(c.RemoteUrl, "image/jpeg", bytes.NewBuffer(inByte))
	if err != nil {
		return
	}
	outByte, err := kmgHttp.ResponseReadAllBody(resp)
	if err != nil {
		return
	}
	outByte, err = kmgCrypto.CompressAndEncryptBytesDecode(kmgRpc_Demo_encryptKey, outByte)
	if err != nil {
		return
	}
	if len(outByte) == 0 {
		return errors.New("len(outByte)==0")
	}
	switch outByte[0] {
	case kmgRpc_Demo_ResponseCodeError:
		return errors.New(string(outByte[1:]))
	case kmgRpc_Demo_ResponseCodeSuccess:
		return json.Unmarshal(outByte[1:], outData)
	default:
		return fmt.Errorf("httpjsonApi protocol error 1 %d", outByte[0])
	}
}

func (s *generateServer_Demo) handleApiV1(inBuf []byte) (outBuf []byte, err error) {
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

	case "PostScoreInt":

		var Info string

		var Err error
		reqData := &struct {
			LbId string

			Score int
		}{}
		Err = json.Unmarshal(b2, reqData)
		if Err != nil {
			return nil, Err
		}

		Info, Err = s.obj.PostScoreInt(reqData.LbId, reqData.Score)
		if Err != nil {
			return nil, Err
		}

		return json.Marshal(struct {
			Info string
		}{

			Info: Info,
		})

	}
	return nil, fmt.Errorf("api %s not found", name)
}

func (c *Client_Demo) PostScoreInt(LbId string, Score int) (Info string, Err error) {
	reqData := &struct {
		LbId string

		Score int
	}{

		LbId: LbId,

		Score: Score,
	}

	respData := &struct {
		Info string
	}{}
	Err = c.sendRequest("PostScoreInt", reqData, &respData)
	return respData.Info, Err

}
