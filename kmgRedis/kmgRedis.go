package kmgRedis

import (
	"fmt"

	"github.com/bronze1man/kmg/encoding/kmgGob"
	"github.com/bronze1man/kmg/encoding/kmgYaml"
	"gopkg.in/redis.v3"
)

// 清空当前db的所有数据,并且设置为传入的数据
// 这个函数主要用于测试
func MustSetDbData(c *redis.Client, data map[string]string) {
	result, err := c.FlushDb().Result()
	if err != nil && result != "OK" {
		panic(fmt.Errorf("[MustSetRedisData] redisClient.FlushDb() fail %s %s", err, result))
	}
	for key, value := range data {
		result, err := c.Set(key, value,0).Result()
		if err != nil && result != "OK" {
			panic(fmt.Errorf("[MustSetRedisData] rc.Set(key,value).Result() fail %s %s", err, result))
		}
	}
}

func MustSetDbDataYaml(c *redis.Client, yaml string) {
	data := map[string]string{}
	err := kmgYaml.Unmarshal([]byte(yaml), &data)
	if err != nil {
		panic(err)
	}
	MustSetDbData(c, data)
	return
}

func MustFlushDb(c *redis.Client) {
	result, err := c.FlushDb().Result()
	if err != nil && result != "OK" {
		panic(fmt.Errorf("[MustFlushDb] redisClient.FlushDb() fail %s %s", err, result))
	}
}

func MustSetDataWithGob(c *redis.Client, key string, data interface{}) {
	b, err := kmgGob.Marshal(data)
	if err != nil {
		panic(err)
	}
	result, err := c.Set(key, string(b),0).Result()
	if err != nil && result != "OK" {
		panic(fmt.Errorf("[MustSetDataWithGob] rc.Set(key,value).Result() fail %s %s", err, result))
	}
}

func MustSetData(c *redis.Client, key string, data string) {
	result, err := c.Set(key, data,0).Result()
	if err != nil && result != "OK" {
		panic(fmt.Errorf("[MustSetData] rc.Set(key,value).Result() fail %s %s", err, result))
	}
}

func MustGetData(c *redis.Client, key string) string {
	result, err := c.Get(key).Result()
	if err != nil {
		panic(fmt.Errorf("[MustGetData] fail %s", err))
	}
	return result
}

func MustGetDataWithGob(c *redis.Client, key string, inData interface{}) {
	result, err := c.Get(key).Result()
	if err != nil {
		panic(fmt.Errorf("[MustGetDataWithGob] fail %s", err))
	}
	err = kmgGob.Unmarshal([]byte(result), inData)
	if err != nil {
		panic(err)
	}
}

func MustGetAllZRange(c *redis.Client, key string) []string {
	result, err := c.ZRange(key, 0, -1).Result()
	if err != nil {
		panic(err)
	}
	return result
}
