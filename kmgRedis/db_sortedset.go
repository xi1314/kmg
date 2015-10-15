package kmgRedis
import (
	"gopkg.in/redis.v3"
)

/*
向跳跃表插入一条数据
向一条不是sortset的数据里面插入,会返回 ErrSortedSetWrongType
网络错误会返回error
*/
func ZAdd(key string, score float64, member string) (err error) {
	err = gClient.ZAdd(key, redis.Z{
		Score:  score,
		Member: member,
	}).Err()
	if err == nil {
		return
	}
	if isRedisErrorWrongType(err) {
		return ErrSortedSetWrongType
	}
	return err
}

func MustZAdd(key string, score float64, member string)  {
	err := ZAdd(key,score,member)
	if err!=nil{
		panic(err)
	}
}

/*
正向读取跳跃表中所有数据
读取一条不是sortset的数据,会返回 ErrSortedSetWrongType
网络错误会返回error
*/
func GetAllScoreAndMemberFromSortedSet(key string) (outList []Z, err error) {
	outList1, err := gClient.ZRangeWithScores(key, 0, -1).Result()
	if err == nil {
		return ZListFromRedisZ(outList1), err
	}
	if isRedisErrorWrongType(err) {
		return nil, ErrSortedSetWrongType
	}
	return nil, err
}

/*
逆向读取跳跃表中所有数据
读取一条不是sortset的数据,会返回 ErrSortedSetWrongType
网络错误会返回error
*/
func GetRevAllScoreAndMemberFromSortedSet(key string) (outList []Z, err error) {
	return ZRevRangeWithScore(key, 0, -1)
}

func ZRevRange(key string, start int, end int) (sList []string, err error) {
	sList, err = gClient.ZRevRange(key, int64(start), int64(end)).Result()
	if err == nil {
		return sList, err
	}
	if isRedisErrorWrongType(err) {
		return nil, ErrSortedSetWrongType
	}
	return nil, err
}

type Z struct{
	Score  float64
	Member string
}

func ZRevRangeWithScore(key string, start int, end int) (outList []Z, err error) {
	outList1, err := gClient.ZRevRangeWithScores(key,int64(start), int64(end)).Result()
	if err == nil {
		return ZListFromRedisZ(outList1), err
	}
	if isRedisErrorWrongType(err) {
		return nil, ErrSortedSetWrongType
	}
	return nil, err
}

func ZListFromRedisZ(list []redis.Z) []Z{
	out:=make([]Z,len(list))
	for i:=range list{
		out[i].Member = list[i].Member.(string)
		out[i].Score = list[i].Score
	}
	return out
}