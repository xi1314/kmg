package demo

import (
	"github.com/bronze1man/kmg/encoding/kmgGob"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"net/http"
)

func ListenAndServe_Demo(addr string, demo *Demo) {
	s := &generateServer_Demo{
		demo: demo,
	}
	err := http.ListenAndServe(addr, s)
	if err != nil {
		panic(err)
	}
}

type generateServer_Demo struct {
	demo *Demo
}

func (s *generateServer_Demo) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.NotFound(w, req)
		panic(`req.Method!="POST"`)
		//kmgLog.Log("InfoServerError", `req.Method!="POST"`, kmgHttp.NewLogStruct(req))
		return
	}
	defer req.Body.Close()
	b1, err := kmgHttp.RequestReadAllBody(req)
	if err != nil {
		panic(err)
		//http.Error(w, "error 1", 400)
		//kmgLog.Log("InfoServerError", err.Error(), kmgHttp.NewLogStruct(req))
		return
	}
	rpcReq := &RpcRequest{}
	kmgGob.MustUnmarshal(b1, rpcReq)
	resp := s.protocol(rpcReq)
	out := kmgGob.MustMarshal(resp)
	w.Write(out)
	return
}

func (s *generateServer_Demo) protocol(req *RpcRequest) (resp RpcResponse) {
	switch req.ApiName {
	case "PostScoreInt":
		if len(req.InData) != 2 {
			resp.Error = "PostScoreInt function parameters number not match"
			return
		}
		var p1 string
		var p2 int
		kmgGob.MustUnmarshal(req.InData[0], &p1)
		kmgGob.MustUnmarshal(req.InData[1], &p2)
		err := s.demo.PostScoreInt(p1, p2)
		if err != nil {
			resp.Error = err.Error()
			return
		}
		return
	case "GetMaxScoreInt":
	case "AutoRegister":
	case "GetFrontUserInfo":
	case "SimpleAdd":
	}
	return
}
