package ajkApi

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/bronze1man/kmg/kmgContext"
	"github.com/bronze1man/kmg/sessionStore"
)

type ApiFuncArgumentError struct {
	Reason  string
	ApiName string
}

func (err *ApiFuncArgumentError) Error() string {
	return fmt.Sprintf("api argument error, reason:%s, name:%s", err.Reason, err.ApiName)
}

type ApiFuncNotFoundError struct {
	Reason  string
	ApiName string
}

func (err *ApiFuncNotFoundError) Error() string {
	return fmt.Sprintf("api function not found, reason:%s, name:%s", err.Reason, err.ApiName)
}

var DefaultApiManager RegisterApiManager = NewApiManager()

//api注册方式
type RegisterApiManager map[string]func(c kmgContext.Context) interface{}

type tKey uint8

const sessionKey tKey = 1

func ContextSetSession(c kmgContext.Context, sess *sessionStore.Session) {
	c.SetValue(sessionKey, sess)
}
func ContextGetSession(c kmgContext.Context) *sessionStore.Session {
	return c.Value(sessionKey).(*sessionStore.Session)
}

/*
 container service + method -> api
 the api name will be "serviceName.methodName"
*/
func NewApiManager() RegisterApiManager {
	return make(RegisterApiManager)
}

func (manager RegisterApiManager) RpcCall(
	session *sessionStore.Session,
	name string,
	caller func(*ApiFuncMeta) error,
) error {
	dotP := strings.LastIndex(name, ".")
	if dotP == -1 {
		return &ApiFuncNotFoundError{Reason: "name not cantain .", ApiName: name}
	}
	c := kmgContext.NewContext()
	ContextSetSession(c, session)

	serviceName := name[:dotP]
	serviceIniter, ok := manager[serviceName]
	if !ok {
		return &ApiFuncNotFoundError{Reason: "service not exist", ApiName: name}
	}
	service := serviceIniter(c)

	serviceType := reflect.TypeOf(service)
	methodName := name[dotP+1:]
	method, ok := serviceType.MethodByName(methodName)
	if ok == false {
		return &ApiFuncNotFoundError{Reason: "method not on service", ApiName: name}
	}
	return caller(&ApiFuncMeta{
		IsMethod:     true,
		Func:         method.Func,
		AttachObject: reflect.ValueOf(service),
		MethodName:   methodName,
	})
}
