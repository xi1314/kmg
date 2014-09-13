package kmgContext

//一个上下文对象,可以保存当前上下文对象的一些内容,也可以修改这些内容
type Context interface {
	Value(key interface{}) interface{}
	SetValue(key interface{}, val interface{})
}

//返回一个新的上下文对象,内部没有锁
func NewContext() Context {
	return make(context)
}

type context map[interface{}]interface{}

func (c context) Value(key interface{}) interface{} {
	return c[key]
}
func (c context) SetValue(key interface{}, val interface{}) {
	c[key] = val
}
