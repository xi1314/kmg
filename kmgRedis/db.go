package kmgRedis

import (
	"fmt"
	"github.com/bronze1man/kmg/encoding/kmgGob"
	"gopkg.in/redis.v3"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var gClient *redis.Client

func DefaultInit() {
	gClient = redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "127.0.0.1:6379",
		DB:      0,
	})
}

func TestInit() {
	gClient = redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "127.0.0.1:6379",
		DB:      1,
	})
}

func InitWithConfig(opt *redis.Options) {
	gClient = redis.NewClient(opt)
}

// 向redis中插入数据,如果已经存在数据会返回 ErrKeyExist
// 如果出现网络错误,会返回 一个网络错误的err
// 没有其他错误的可能性了.
func Insert(key string, value string) (err error) {
	success, err := gClient.SetNX(key, value,0).Result()
	if err != nil {
		return
	}
	if !success {
		return ErrKeyExist
	}
	return nil
}

// 向redis中插入数据,如果已经存在数据会返回 ErrKeyExist
func InsertGob(key string, obj interface{}) (err error) {
	b, err := kmgGob.Marshal(obj)
	if err != nil {
		return
	}
	return Insert(key, string(b))
}

func MustInsert(key string, value string) {
	err := Insert(key, value)
	if err != nil {
		panic(err)
	}
}

func MustInsertGob(key string, obj interface{}) {
	err := InsertGob(key, obj)
	if err != nil {
		panic(err)
	}
}

// 数据不存在会返回 ErrKeyNotExist
// 网络错误会返回 error
func Update(key string, value string) (err error) {
	cmd := redis.NewStatusCmd("SET", key, value, "XX")
	gClient.Process(cmd)
	_, err = cmd.Result()
	if err == redis.Nil {
		return ErrKeyNotExist
	}
	return err
}

// key,不存在会insert,存在会update
// 网络错误会返回 error
// 注意redis类型错误时,也不会报错,只会把这个key设置成正确的value
func Set(key string, value string) (err error) {
	return gClient.Set(key, value,0).Err()
}

func MustSet(key string, value string) {
	err := Set(key, value)
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

// 从redis的kvdb中获取一个key
// 注意 value有可能是 "" 这个和数据不存在是两种情况.
// 如果数据不存在,会返回ErrKeyNotExist
// value在redis里面不是string类型,会返回 ErrStringWrongType
// 网络错误也会返回 error
func Get(key string) (value string, err error) {
	value, err = gClient.Get(key).Result()
	if err==nil{
		return value,nil
	}
	if isRedisErrorWrongType(err){
		return "",ErrStringWrongType
	}
	if err == redis.Nil {
		return "", ErrKeyNotExist
	}
	return "", err
}

func MustGet(key string) (value string) {
	value, err := Get(key)
	if err != nil {
		panic(err)
	}
	return value
}

// 从redis的kvdb中获取一个key
// 将这个key转换成int
// 无法转换成int,会panic
// key不存在,返回0
// 网络错误会panic
func MustGetIntIgnoreNotExist(key string) (valueI int){
	value,err := Get(key)
	if err==ErrKeyNotExist{
		return 0
	}
	if err!=nil{
		panic(err)
	}
	valueI,err = strconv.Atoi(value)
	if err!=nil{
		panic(err)
	}
	return valueI
}

// 从redis的kvdb中获取一个key
// 将这个key转换成float
// 无法转换成float,会panic
// key不存在,返回0
// 网络错误会panic
func MustGetFloatIgnoreNotExist(key string) (float64){
	value,err := Get(key)
	if err==ErrKeyNotExist{
		return 0
	}
	if err!=nil{
		panic(err)
	}
	valueF,err := strconv.ParseFloat(value,64)
	if err!=nil{
		panic(err)
	}
	return valueF
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

// 只有网络问题会返回error
func Del(key string) (err error) {
	return gClient.Del(key).Err()
}

// 只有网络问题会返回error
func FlushDbV2() (err error) {
	_, err = gClient.FlushDb().Result()
	return err
}
func MustFlushDbV2() {
	err := FlushDbV2()
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

// 使用 redis的表达式搜索key,返回搜索到的key的列表
// 只有网络问题会返回error
func Keys(searchKey string) (keyList []string, err error) {
	return gClient.Keys(searchKey).Result()
}

// 某个key是否存在
// 只有网络问题会返回error
func Exists(key string) (exist bool, err error) {
	return gClient.Exists(key).Result()
}

/*
Insert all the specified values at the tail of the list stored at key. If key does not exist, it is created as empty list before performing the push operation. When key holds a value that is not a list, an error is returned.
更改的key存在,会向这个数组类型的key,右边加入一个元素.
更改的key不存在,会创建一个,并且写入第一个值.
更改的key的类型不正确会返回 ErrListWrongType
网络错误会返回error
*/
func RPush(key string, value string) (err error) {
	err = gClient.RPush(key, value).Err()
	if err == nil {
		return nil
	}
	if isRedisErrorWrongType(err) {
		return ErrListWrongType
	}
	return err
}

/*
返回一个redis数组里面所有的值.
查询的key存在,并且类型正确,返回列表里面的数据
查询的key不存在,返回空数组 TODO 好用?
查询的key类型错误,返回 ErrListWrongType
网络错误会返回error
*/
func GetAllValueInList(key string) (out []string, err error) {
	out, err = gClient.LRange(key, 0, -1).Result()
	if err == nil {
		return out, nil
	}
	if isRedisErrorWrongType(err) {
		return nil, ErrListWrongType
	}
	return nil, err
}

/*
一次操作,批量从redis里面返回大量key的值.
没有传入数据,不报错
如果查询的key全部存在,返回数据.
如果存在某一个key不存在,或者类型错误,返回 ErrKeyNotExist ,value里面什么也没有 (和redis命令不一致)
网络错误会返回error
*/
func MGetNotExistCheck(keyList []string) (value []string, err error) {
	if len(keyList) == 0 {
		return nil, nil
	}
	outList, err := gClient.MGet(keyList...).Result()
	if err != nil {
		return nil, err
	}
	value = make([]string, len(outList))
	for i, stringI := range outList {
		s, ok := stringI.(string)
		if !ok {
			return nil, ErrKeyNotExist
		}
		value[i] = s
	}
	return value, nil
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
	return mgetNotExistCheckGobUnmarshal(outList, reflect.ValueOf(obj))
}

func mgetNotExistCheckGobUnmarshal(outList []interface{}, obj reflect.Value) (err error) {
	switch obj.Kind() {
	case reflect.Ptr:
		return mgetNotExistCheckGobUnmarshal(outList, obj.Elem())
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

/*
带超时的设置一条数据
网络错误会返回error
*/
func SetEx(key string, dur time.Duration, value string) (err error) {
	return gClient.Set(key, value,dur).Err()
}

func isRedisErrorWrongType(err error) bool {
	return strings.Contains(err.Error(), "WRONGTYPE")
}

/*
改key的名字
key不存在       ErrKeyNotExist
key和newKey一样 ErrRenameSameName
newKey存在      ErrKeyExist
网络错误会返回error
*/
func RenameNx(key string, newKey string) (err error) {
	retB, err := gClient.RenameNX(key, newKey).Result()
	if err == nil {
		if retB == false {
			return ErrKeyExist
		}
		return
	}
	errS := err.Error()
	if strings.Contains(errS, "ERR source and destination objects are the same") {
		return ErrRenameSameName
	}
	if strings.Contains(errS, "ERR no such key") {
		return ErrKeyNotExist
	}
	return err
}

/*
给某一个redis的key加一个整数
key不存在,会先把这个key变成0,然后再进行增加
key不能被解析成整数,会返回 ErrValueNotIntFormatOrOutOfRange
value不是string类型,会返回 ErrStringWrongType
网络错误会返回error
 */
func IncrBy(key string,num int64) (err error){
	err=gClient.IncrBy(key,num).Err()
	if err!=nil{
		if isRedisErrorWrongType(err){
			return ErrStringWrongType
		}
		if strings.Contains(err.Error(),"ERR value is not an integer or out of range"){
			return ErrValueNotIntFormatOrOutOfRange
		}
		return
	}
	return nil
}

/*
给某一个redis的key加一个浮点
不要传入大于1e200的浮点,会挂掉. TODO 解决这个问题?
key不存在,会先把这个key变成0,然后再进行增加
key不能被解析成整数,会返回 ErrValueNotFloatFormatOrOutOfRange
value不是string类型,会返回 ErrStringWrongType
网络错误会返回error
 */
func IncrByFloat(key string,num float64) (err error){
	err=gClient.IncrByFloat(key,num).Err()
	if err!=nil{
		if isRedisErrorWrongType(err){
			return ErrStringWrongType
		}
		if strings.Contains(err.Error(),"ERR value is not a valid float"){
			return ErrValueNotFloatFormat
		}
		return
	}
	return nil
}

