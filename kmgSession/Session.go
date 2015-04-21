package kmgSession

import (
	"github.com/bronze1man/kmg/kmgRand"
	"net/http"
	"sync"
)

type Session struct {
	Id   string
	Data map[string]string
	lock sync.Mutex
}

var SessionIdName string = "KmgSessionId"

var Store = map[string]*Session{}

//通过 Session Id 获取 Session，如果没获取到，产生一个新的返回
func GetSessionById(id string) *Session {
	if session, ok := Store[id]; ok {
		return session
	} else {
		session = &Session{
			Id:   kmgRand.MustCryptoRandToAlphaNum(26),
			Data: map[string]string{},
		}
		Store[session.Id] = session
		return session
	}
}

func GetSession(w http.ResponseWriter, req *http.Request) *Session {
	var id string = ""
	cookie, _ := req.Cookie(SessionIdName)
	if cookie != nil {
		id = cookie.Value
	} else {
		cookie = &http.Cookie{
			Name: SessionIdName,
		}
	}
	session := GetSessionById(id)
	cookie.Value = session.Id
	http.SetCookie(w, cookie)
	return session
}

func (s *Session) Set(k string, v string) {
	s.Data[k] = v
}

func (s *Session) Get(k string) string {
	if v, ok := s.Data[k]; ok {
		return v
	} else {
		return ""
	}
}
