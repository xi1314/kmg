package kmgRedis

import (
	"github.com/bronze1man/kmg/kmgStrconv"
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

func MustZAdd(key string, score float64, member string) {
	err := ZAdd(key, score, member)
	if err != nil {
		panic(err)
	}
}

/*
使用member获取Score
读取一条不存在的key,会返回 ErrKeyNotExist
key存在,但是Member,不存在会返回 ErrKeyNotExist
读取一条不是sortset的数据,会返回 ErrSortedSetWrongType
网络错误会返回error
*/
func ZScore(key string,member string) (f float64,err error){
	f,err = gClient.ZScore(key,member).Result()
	if err==nil{
		return f,nil
	}
	if err==redis.Nil{
		return 0,ErrKeyNotExist
	}
	if isRedisErrorWrongType(err) {
		return 0,ErrSortedSetWrongType
	}
	return 0,err
}

func MustZScore(key string,member string) (f float64){
	f,err  := ZScore(key,member)
	if err!=nil{
		panic(err)
	}
	return f
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

func MustGetAllMemberFromSortedSet(key string) (sList []string){
	sList,err := ZRange(key,0,-1)
	if err!=nil{
		panic(err)
	}
	return sList
}

/*
逆向读取跳跃表中所有数据
读取一条不是sortset的数据,会返回 ErrSortedSetWrongType
网络错误会返回error
*/
func GetRevAllScoreAndMemberFromSortedSet(key string) (outList []Z, err error) {
	return ZRevRangeWithScore(key, 0, -1)
}

func ZRange(key string, start int, end int) (sList []string, err error) {
	sList, err = gClient.ZRange(key, int64(start), int64(end)).Result()
	if err == nil {
		return sList, err
	}
	if isRedisErrorWrongType(err) {
		return nil, ErrSortedSetWrongType
	}
	return nil, err
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

type Z struct {
	Score  float64
	Member string
}

func ZRevRangeWithScore(key string, start int, end int) (outList []Z, err error) {
	outList1, err := gClient.ZRevRangeWithScores(key, int64(start), int64(end)).Result()
	if err == nil {
		return ZListFromRedisZ(outList1), err
	}
	if isRedisErrorWrongType(err) {
		return nil, ErrSortedSetWrongType
	}
	return nil, err
}

func ZRangeByScoreWithScore(key string, min float64,max float64) (zList []Z,err error){
	zList1, err := gClient.ZRangeByScoreWithScores(key, redis.ZRangeByScore{
		Min: kmgStrconv.FormatFloat(min),
		Max: kmgStrconv.FormatFloat(max),
	}).Result()
	if err == nil {
		return ZListFromRedisZ(zList1), err
	}
	if isRedisErrorWrongType(err) {
		return nil, ErrSortedSetWrongType
	}
	return nil, err
}

func ZRangeByScore(key string, min float64, max float64) (sList []string, err error) {
	sList, err = gClient.ZRangeByScore(key, redis.ZRangeByScore{
		Min: kmgStrconv.FormatFloat(min),
		Max: kmgStrconv.FormatFloat(max),
	}).Result()
	if err == nil {
		return sList, err
	}
	if isRedisErrorWrongType(err) {
		return nil, ErrSortedSetWrongType
	}
	return nil, err
}

func ZRevRangeByScore(key string, min float64, max float64) (sList []string, err error) {
	sList, err = gClient.ZRevRangeByScore(key, redis.ZRangeByScore{
		Min: kmgStrconv.FormatFloat(min),
		Max: kmgStrconv.FormatFloat(max),
	}).Result()
	if err == nil {
		return sList, err
	}
	if isRedisErrorWrongType(err) {
		return nil, ErrSortedSetWrongType
	}
	return nil, err
}

func ZListFromRedisZ(list []redis.Z) []Z {
	out := make([]Z, len(list))
	for i := range list {
		out[i].Member = list[i].Member.(string)
		out[i].Score = list[i].Score
	}
	return out
}

func ZScanCallback(key string, cb func(member string) error) (err error) {
	pos := 0
	for {
		memberList, err := ZRange(key, pos, pos+scanSize-1)
		if err != nil {
			return err
		}
		for _, member := range memberList {
			err = cb(member)
			if err != nil {
				return err
			}
		}
		if len(memberList) < scanSize {
			// 如果没有数据表示扫描完毕了.
			return nil
		}
		pos += scanSize // 继续扫描下一组数据.
	}
}

func MustGetSortedSetSize(key string) int{
	num,err:= gClient.ZCard(key).Result()
	if err!=nil{
		panic(err)
	}
	return int(num)
}

func MustZRemRangeByScore(key string,min float64,max float64){
	err:=gClient.ZRemRangeByScore(key,kmgStrconv.FormatFloat(min),kmgStrconv.FormatFloat(max)).Err()
	if err!=nil{
		panic(err)
	}
}