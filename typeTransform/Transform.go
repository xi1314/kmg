package typeTransform

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgReflect"
	"github.com/bronze1man/kmg/kmgStrconv"
	"github.com/bronze1man/kmg/kmgTime"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

/*
try best transform one type to another type
special case:
"" => 0
"" => 0.0
*/
func Transform(in interface{}, out interface{}) (err error) {
	return DefaultTransformer.Transform(in, out)
}

func MustTransform(in interface{}, out interface{}) {
	err := DefaultTransformer.Transform(in, out)
	if err != nil {
		panic(err)
	}
}

func MustTransformToMap(in interface{}) (m map[string]string) {
	m = map[string]string{}
	err := DefaultTransformer.Transform(in, &m)
	if err != nil {
		panic(err)
	}
	return m
}

func MapToMap(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	out.Set(reflect.MakeMap(out.Type()))
	for _, key := range in.MapKeys() {
		oKey := reflect.New(out.Type().Key()).Elem()
		oVal := reflect.New(out.Type().Elem()).Elem()
		err = t.Tran(key, oKey)
		if err != nil {
			return
		}
		val := in.MapIndex(key)
		err = t.Tran(val, oVal)
		if err != nil {
			return
		}
		out.SetMapIndex(oKey, oVal)
	}
	return
}

func StringToString(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	out.SetString(in.String())
	return nil
}

func NewStringToTimeFunc(location *time.Location) TransformerFunc {
	return func(traner Transformer, in reflect.Value, out reflect.Value) (err error) {
		var t time.Time
		t, err = kmgTime.ParseAutoInLocation(in.String(), location)
		if err != nil {
			return
		}
		out.Set(reflect.ValueOf(t))
		return
	}
}
func StringToTime(traner Transformer, in reflect.Value, out reflect.Value) (err error) {
	var t time.Time
	t, err = kmgTime.ParseAutoInLocal(in.String())
	if err != nil {
		return
	}
	out.Set(reflect.ValueOf(t))
	return
}

func TimeToString(traner Transformer, in reflect.Value, out reflect.Value) (err error) {
	t := in.Interface().(time.Time)
	out.SetString(t.In(kmgTime.DefaultTimeZone).Format(kmgTime.FormatMysql))
	return nil
}

func PtrToPtr(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	return t.Tran(in.Elem(), out.Elem())
}

//假设map的key类型是string,值类型不限
func MapToStruct(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	oKey := reflect.New(reflect.TypeOf("")).Elem()
	out.Set(reflect.New(out.Type()).Elem())
	fieldNameMap := map[string]bool{}
	for _, key := range in.MapKeys() {
		err = t.Tran(key, oKey)
		if err != nil {
			return
		}
		sKey := oKey.String()
		fieldNameMap[sKey] = true
		oVal := out.FieldByName(sKey)
		if !oVal.IsValid() {
			continue
		}
		val := in.MapIndex(key)
		err = t.Tran(val, oVal)
		if err != nil {
			return
		}
	}
	return
}

//假设map的key类型是string,值类型不限
func StructToMap(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	oValType := out.Type().Elem()
	oKey := reflect.New(reflect.TypeOf("")).Elem()
	if out.IsNil() {
		out.Set(reflect.MakeMap(out.Type())) //新建map不能使用reflect.New
	}

	fieldMap := kmgReflect.StructGetAllFieldMap(in.Type())
	for key, field := range fieldMap {
		if field.PkgPath != "" {
			//忽略没有导出的字段
			continue
		}
		iVal := in.FieldByName(key)
		oVal := reflect.New(oValType).Elem()
		err = t.Tran(iVal, oVal)
		if err != nil {
			return
		}
		oKey.SetString(key)
		out.SetMapIndex(oKey, oVal)
	}
	return nil
}

func SliceToSlice(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	len := in.Len()
	out.Set(reflect.MakeSlice(out.Type(), len, len))
	for i := 0; i < len; i++ {
		val := in.Index(i)
		err = t.Tran(val, out.Index(i))
		if err != nil {
			return
		}
	}
	return
}

// "" => 0
func StringToInt(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	inS := in.String()
	inS = strings.TrimSpace(inS)
	if inS == "" {
		out.SetInt(int64(0))
		return nil
	}
	i, err := strconv.ParseInt(inS, 10, out.Type().Bits())
	if err != nil {
		return
	}
	out.SetInt(i)
	return
}

// "" => 0
func StringToUint(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	inS := in.String()
	inS = strings.TrimSpace(inS)
	if inS == "" {
		out.SetUint(uint64(0))
		return nil
	}
	i, err := strconv.ParseUint(inS, 10, out.Type().Bits())
	if err != nil {
		return
	}
	out.SetUint(i)
	return
}

// "" => 0.0
func StringToFloat(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	inS := in.String()
	if inS == "" {
		out.SetFloat(0.0)
		return nil
	}
	i, err := strconv.ParseFloat(inS, out.Type().Bits())
	if err != nil {
		return
	}
	out.SetFloat(i)
	return
}

// "" => false
func StringToBool(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	inS := in.String()
	if inS == "" {
		out.SetBool(false)
		return nil
	}
	i, err := strconv.ParseBool(inS)
	if err != nil {
		return
	}
	out.SetBool(i)
	return
}

func IntToInt(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	out.SetInt(in.Int())
	return nil
}

func FloatToInt(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	outf1 := in.Float()
	if math.Floor(outf1) != outf1 {
		return fmt.Errorf("[typeTransform.tran] it seems to lose some accuracy trying to convert from float to int,float:%f", outf1)
	}
	out.SetInt(int64(outf1))
	return
}

func FloatToFloat(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	out.SetFloat(in.Float())
	return
}

func FloatToString(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	f := in.Float()
	fs := kmgStrconv.FormatFloat(f)
	out.SetString(fs)
	return
}

func NonePtrToPtr(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	if out.IsNil() {
		out.Set(reflect.New(out.Type().Elem()))
	}
	return t.Tran(in, out.Elem())
}
func InterfaceToNoneInterface(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	return t.Tran(in.Elem(), out)
}
func NoneInterfaceToInterface(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	//满足interface的情况
	if in.Type().Implements(out.Type()) {
		out.Set(in)
		return
	}
	return t.Tran(in, out.Elem())
}

func IntToString(t Transformer, in reflect.Value, out reflect.Value) (err error) {
	s := strconv.FormatInt(in.Int(), 10)
	out.SetString(s)
	return nil
}
