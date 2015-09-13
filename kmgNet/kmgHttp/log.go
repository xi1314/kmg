package kmgHttp

import (
	"net/http"
)

type LogStruct struct {
	Method     string
	Url        string
	RemoteAddr string
	UA         string `json:",omitempty"`
	Refer      string `json:",omitempty"`
	Host       string `json:",omitempty"`
}

func NewLogStruct(req *http.Request) *LogStruct {
	return &LogStruct{
		Method:     req.Method,
		Url:        req.URL.String(),
		RemoteAddr: req.RemoteAddr,
		UA:         req.UserAgent(),
		Refer:      req.Referer(),
		Host:       req.Host,
	}
}
