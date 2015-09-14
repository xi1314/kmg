package kmgLog

import (
	"encoding/json"
	"time"
)

type LogWriter func(r LogRow)

// 这个LogRow仅仅是用来序列化的,拿去反序列化不好用,不靠谱
type LogRow struct {
	Cat  string
	Time time.Time // 让json库自己折腾去.
	// Data []json.RawMessage //反序列化可以用这个定义,但是此处只有序列化用定义
	Data []interface{}
}

func (r LogRow) Marshal() (b []byte, err error) {
	return json.Marshal(r)
}

/*
func (r LogRow) UnmarshalData(index int, obj interface{}) (err error) {
	return json.Unmarshal(r.Data[index], obj)
}
func (r LogRow) MustUnmarshalData(index int, obj interface{}) {
	err := json.Unmarshal(r.Data[index], obj)
	if err != nil {
		panic(err)
	}
}
*/