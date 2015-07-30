package kmgSession

import (
	"net/http"
	"sync"

	"github.com/bronze1man/kmg/kmgRand"
)

type Session struct {
	Id   string
	Data map[string]string
	lock sync.RWMutex //此处不用锁会违反map的线性化写入的规则.
}

var SessionIdName string = "KmgSessionId"

var MemoryStore = Store{
	jar: map[string]*Session{},
}

type Store struct {
	lock sync.Mutex
	jar  map[string]*Session
}

func (store *Store) Get(sessionId string) (session *Session, ok bool) {
	store.lock.Lock()
	defer store.lock.Unlock()
	session, ok = store.jar[sessionId]
	return session, ok
}

func (store *Store) Set(sessionId string, session *Session) {
	store.lock.Lock()
	defer store.lock.Unlock()
	store.jar[sessionId] = session
}

//通过 Session Id 获取 Session，如果没获取到，产生一个新的返回
func GetSessionById(sessionId string) *Session {
	if session, ok := MemoryStore.Get(sessionId); ok {
		return session
	} else {
		session = &Session{
			Id:   kmgRand.MustCryptoRandToAlphaNum(26),
			Data: map[string]string{},
		}
		MemoryStore.Set(session.Id, session)
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

func (s *Session) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Data = map[string]string{}
}

func (s *Session) Set(k string, v string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Data[k] = v
}

func (s *Session) Get(k string) string {
	if v, ok := s.Data[k]; ok {
		return v
	} else {
		return ""
	}
}
