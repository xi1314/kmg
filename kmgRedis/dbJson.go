package kmgRedis
import (
	"encoding/json"
	"github.com/bronze1man/kmg/encoding/kmgJson"
)

func GetJson(key string,obj interface{}) (err error){
	value, err := Get(key)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(value), obj)
	if err != nil {
		return err
	}
	return nil
}

func MustInsertJson(key string,obj interface{}) {
	b := kmgJson.MustMarshal(obj)
	MustInsert(key, string(b))
}

func MustUpdateJson(key string,obj interface{}) {
	b := kmgJson.MustMarshal(obj)
	MustUpdate(key, string(b))
}