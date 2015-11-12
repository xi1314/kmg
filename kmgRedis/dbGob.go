package kmgRedis
import (
	"github.com/bronze1man/kmg/encoding/kmgGob"
	"github.com/bronze1man/kmg/typeTransform"
	"reflect"
	"fmt"
)

// 向redis中插入数据,如果已经存在数据会返回 ErrKeyExist
func InsertGob(key string, obj interface{}) (err error) {
	b, err := kmgGob.Marshal(obj)
	if err != nil {
		return
	}
	return Insert(key, string(b))
}

func MustInsertGob(key string, obj interface{}) {
	err := InsertGob(key, obj)
	if err != nil {
		panic(err)
	}
}

func SetGob(key string, obj interface{}) (err error) {
	b, err := kmgGob.Marshal(obj)
	if err != nil {
		return
	}
	return Set(key, string(b))
}

func MustSetGob(key string, obj interface{}) {
	err := SetGob(key, obj)
	if err != nil {
		panic(err)
	}
}

func GetGobIgnoreNotExist(key string, obj interface{}) (err error) {
	err = GetGob(key, obj)
	if err == ErrKeyNotExist {
		return nil
	}
	return err
}

func MustGetGob(key string, obj interface{}) {
	err := GetGob(key, obj)
	if err != nil {
		panic(err)
	}
}

// 向redis中更新数据,如果不存在数据,会返回 ErrKeyNotExist
func UpdateGob(key string, obj interface{}) (err error) {
	b, err := kmgGob.Marshal(obj)
	if err != nil {
		return
	}
	return Update(key, string(b))
}

// 如果数据不存在,会返回ErrKeyNotExist
// 序列化错误,会返回 error
// 网络错误也会返回 error
func GetGob(key string, obj interface{}) (err error) {
	value, err := Get(key)
	if err != nil {
		return err
	}
	err = kmgGob.Unmarshal([]byte(value), obj)
	if err != nil {
		return err
	}
	return nil
}

func RPushGob(key string, value interface{}) (err error) {
	b, err := kmgGob.Marshal(value)
	if err != nil {
		return err
	}
	return RPush(key, string(b))
}

func LRangeAllGob(key string, list interface{}) (err error) {
	sList, err := LRangeAll(key)
	iList := []interface{}{}
	typeTransform.MustTransform(&sList,&iList)
	if err != nil {
		return err
	}
	return ListGobUnmarshalNotExistCheck(iList, reflect.ValueOf(list))
}

/*
带超时的设置一条数据
没有传入数据,不报错,不修改obj
网络错误会返回error
*/
func MGetNotExistCheckGob(keyList []string, obj interface{}) (err error) {
	if len(keyList) == 0 {
		return nil
	}
	outList, err := gClient.MGet(keyList...).Result()
	if err != nil {
		return err
	}
	return ListGobUnmarshalNotExistCheck(outList, reflect.ValueOf(obj))
}

func ListGobUnmarshalNotExistCheck(outList []interface{}, obj reflect.Value) (err error) {
	switch obj.Kind() {
	case reflect.Ptr:
		return ListGobUnmarshalNotExistCheck(outList, obj.Elem())
	case reflect.Slice:
		newSlice := reflect.MakeSlice(obj.Type(), len(outList), len(outList))
		elemType := obj.Type().Elem()
		for i, stringI := range outList {
			s, ok := stringI.(string)
			if !ok {
				return ErrKeyNotExist
			}
			thisValue := newSlice.Index(i)
			thisElem := reflect.New(elemType)
			err = kmgGob.Unmarshal([]byte(s), thisElem.Interface())
			if err != nil {
				return err
			}
			thisValue.Set(thisElem.Elem())
		}
		obj.Set(newSlice)
		return nil
	default:
		return fmt.Errorf("[mgetNotExistCheckGobUnmarshal] Unmarshal unexpect Kind %s", obj.Kind().String())
	}
}