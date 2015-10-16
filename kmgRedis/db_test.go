package kmgRedis

import (
	"github.com/bronze1man/kmg/kmgTest"
	"sort"
	"testing"
	"time"
	"github.com/bronze1man/kmg/kmgStrings"
	"strconv"
	"fmt"
)

func init() {
	TestInit()
}

func TestRedisKvdb(ot *testing.T) {
	MustFlushDbV2()
	err := Set("test_1", "abc")
	kmgTest.Equal(err, nil)
	v, err := Get("test_1")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(v, "abc")

	_, err = Get("test_2")
	kmgTest.Equal(err, ErrKeyNotExist)

	err = Del("test_1")
	kmgTest.Equal(err, nil)

	_, err = Get("test_2")
	kmgTest.Equal(err, ErrKeyNotExist)

	err = Insert("test_3", "abcd")
	kmgTest.Equal(err, nil)

	err = Insert("test_3", "abcde")
	kmgTest.Equal(err, ErrKeyExist)

	kmgTest.Equal(MustGet("test_3"), "abcd")

	err = Set("test_3", "abcdef")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(MustGet("test_3"), "abcdef")

	err = Update("test_3", "abcdef")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(MustGet("test_3"), "abcdef")

	err = Update("test_4", "abcdefg")
	kmgTest.Equal(err, ErrKeyNotExist)
	kmgTest.Equal(MustGet("test_3"), "abcdef")

	simpleData := []string{"1", "2"}
	err = InsertGob("test_4", simpleData)
	kmgTest.Equal(err, nil)

	err = SetGob("test_4", simpleData)
	kmgTest.Equal(err, nil)

	getSimpleData := []string{}
	err = GetGob("test_4", &getSimpleData)
	kmgTest.Equal(err, nil)
	kmgTest.Equal(simpleData, getSimpleData)

	simpleData = []string{"1", "2", "3"}
	err = UpdateGob("test_4", simpleData)

	getSimpleData = []string{}
	err = GetGob("test_4", &getSimpleData)
	kmgTest.Equal(err, nil)
	kmgTest.Equal(simpleData, getSimpleData)

	var getSimpleData1 []string
	err = GetGob("test_4", &getSimpleData1)
	kmgTest.Equal(err, nil)
	kmgTest.Equal(simpleData, getSimpleData1)

	var getSimpleData2 *[]string
	err = GetGob("test_4", getSimpleData2)
	kmgTest.Ok(err != nil)

	err = GetGob("test_4", &getSimpleData2)
	kmgTest.Equal(err, nil)
	kmgTest.Equal(simpleData, *getSimpleData2)

	keyList, err := Keys("test_*")
	kmgTest.Equal(err, nil)
	sort.Strings(keyList)
	kmgTest.Equal(keyList, []string{"test_3", "test_4"})

	exist, err := Exists("test_4")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(exist, true)

	exist, err = Exists("test_5")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(exist, false)

	MustInsert("test_5", "")
	v, err = Get("test_5")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(v, "")

	err = SetEx("test_6", time.Second, "abc")
	kmgTest.Equal(err, nil)

	kmgTest.Equal(MustGet("test_6"), "abc")

	time.Sleep(time.Second * 2)
	exist, err = Exists("test_6")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(exist, false)
}

func TestGet(ot *testing.T){
	MustFlushDbV2()
	MustInsert("test_1", "abc")

	v, err := Get("test_1")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(v, "abc")

	_, err = Get("test_2")
	kmgTest.Equal(err, ErrKeyNotExist)

	MustZAdd("test_2",0,"abc")
	_,err = Get("test_2")
	kmgTest.Equal(err, ErrStringWrongType)
}

func TestSet(ot *testing.T){
	MustFlushDbV2()
	MustInsert("test_1", "abc")

	err := Set("test_1","abcd")
	kmgTest.Equal(err, nil)

	err = Set("test_2","abcde")
	kmgTest.Equal(err, nil)

	MustZAdd("test_3",0,"abc")
	err = Set("test_3","abcdefg")
	kmgTest.Equal(err, nil)

	kmgTest.Equal(MustGet("test_3"),"abcdefg")
}

func TestMGetSlice(ot *testing.T) {
	MustFlushDbV2()
	MustInsert("test_1", "abc")
	MustInsert("test_2", "abcd")
	outList, err := MGetNotExistCheck([]string{"test_1", "test_2"})
	kmgTest.Equal(err, nil)
	kmgTest.Equal(outList, []string{"abc", "abcd"})

	outList, err = MGetNotExistCheck([]string{"test_1", "test_2", "test_3"})
	kmgTest.Equal(err, ErrKeyNotExist)
	kmgTest.Equal(outList, nil)
}

func TestMGetNotExistCheckGob(ot *testing.T) {
	MustFlushDbV2()
	MustInsertGob("test_1", []string{"test_1", "test_3"})
	MustInsertGob("test_2", []string{"test_2"})

	var output [][]string
	err := MGetNotExistCheckGob([]string{"test_1", "test_2"}, &output)
	kmgTest.Equal(err, nil)
	kmgTest.Equal(len(output), 2)
	kmgTest.Equal(output[0], []string{"test_1", "test_3"})
	kmgTest.Equal(output[1], []string{"test_2"})

	dataList, err := MGetNotExistCheck(nil)
	kmgTest.Equal(err, nil)
	kmgTest.Equal(len(dataList), 0)

	output = [][]string{}
	err = MGetNotExistCheckGob(nil, output)
	kmgTest.Equal(err, nil)
	kmgTest.Equal(len(output), 0)
}

func TestRedisList(ot *testing.T) {
	MustFlushDbV2()
	MustInsert("test_4", "abc")

	err := RPush("test_4", "abc")
	kmgTest.Equal(err, ErrListWrongType)

	err = RPush("test_5", "abc")
	kmgTest.Equal(err, nil)

	err = RPush("test_5", "abcd")
	kmgTest.Equal(err, nil)

	list, err := GetAllValueInList("test_5")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(list, []string{"abc", "abcd"})

	_, err = GetAllValueInList("test_4")
	kmgTest.Equal(err, ErrListWrongType)

	list, err = GetAllValueInList("test_6")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(len(list), 0)
}

func TestRedisSortedSet(ot *testing.T) {
	MustFlushDbV2()
	MustInsert("test_4", "abc")

	err := ZAdd("test_4", 0, "abc")
	kmgTest.Equal(err, ErrSortedSetWrongType)

	err = ZAdd("test_1", 0, "abc")
	kmgTest.Equal(err, nil)

	err = ZAdd("test_1", -1, "abcd")
	kmgTest.Equal(err, nil)

	MustZAdd("test_1",-1,"abcd")

	zlist, err := GetAllScoreAndMemberFromSortedSet("test_1")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(zlist, []Z{
		{Score: -1, Member: "abcd"},
		{Score: 0, Member: "abc"},
	})

	_, err = GetAllScoreAndMemberFromSortedSet("test_4")
	kmgTest.Equal(err, ErrSortedSetWrongType)

	zlist, err = GetAllScoreAndMemberFromSortedSet("test_5")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(len(zlist), 0)

	zlist, err = GetRevAllScoreAndMemberFromSortedSet("test_1")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(zlist, []Z{
		{Score: 0, Member: "abc"},
		{Score: -1, Member: "abcd"},
	})

	sList, err := ZRevRange("test_1", 0, 0)
	kmgTest.Equal(err, nil)
	kmgTest.Equal(sList, []string{"abc"})

	zlist, err = ZRevRangeWithScore("test_1", 0, 0)
	kmgTest.Equal(err, nil)
	kmgTest.Equal(zlist, []Z{
		{Score: 0, Member: "abc"},
	})

	sList,err = ZRevRangeByScore("test_1",-0.5,-2)
	kmgTest.Equal(err,nil)
	kmgTest.Equal(sList, []string{"abcd"})

	sList,err = ZRevRangeByScore("test_5",-0.5,-2)
	kmgTest.Equal(err, nil)
	kmgTest.Equal(len(sList), 0)

	sList,err = ZRevRangeByScore("test_4",-0.5,-2)
	kmgTest.Equal(err, ErrSortedSetWrongType)
}

func TestRedisRename(ot *testing.T) {
	MustFlushDbV2()
	MustInsert("test_1", "abc")
	MustInsert("test_4", "abc")

	// 正常情况
	err := RenameNx("test_1", "test_2")
	kmgTest.Equal(err, nil)

	// key和newKey一样
	err = RenameNx("test_2", "test_2")
	kmgTest.Equal(err, ErrRenameSameName)

	//key不存在
	err = RenameNx("test_1", "test_3")
	kmgTest.Equal(err, ErrKeyNotExist)

	// newKey存在
	err = RenameNx("test_2", "test_4")
	kmgTest.Equal(err, ErrKeyExist)
}

func TestMustGetIntDefault0(ot *testing.T){
	MustFlushDbV2()
	MustInsert("test_1", "abc")
	MustInsert("test_2", "2")

	kmgTest.AssertPanic(func(){
		MustGetIntIgnoreNotExist("test_1")
	})
	outI:=MustGetIntIgnoreNotExist("test_2")
	kmgTest.Equal(outI,2)

	outI=MustGetIntIgnoreNotExist("test_3")
	kmgTest.Equal(outI,0)
}

func TestMustGetFloatDefault0(ot *testing.T){
	MustFlushDbV2()
	MustInsert("test_1", "abc")
	MustInsert("test_2", "2.1")

	kmgTest.AssertPanic(func(){
		MustGetFloatIgnoreNotExist("test_1")
	})
	outI:=MustGetFloatIgnoreNotExist("test_2")
	kmgTest.Equal(outI,2.1)

	outI=MustGetFloatIgnoreNotExist("test_3")
	kmgTest.Equal(outI,0.0)
}

func TestIncrBy(ot *testing.T){
	MustFlushDbV2()
	MustInsert("test_2","abc")
	MustZAdd("test_3",0,"abc")

	err:=IncrBy("test_1",2)
	kmgTest.Equal(err,nil)
	kmgTest.Equal(MustGetIntIgnoreNotExist("test_1"),2)

	err = IncrBy("test_2",3)
	kmgTest.Equal(err,ErrValueNotIntFormatOrOutOfRange)

	err = IncrBy("test_3",4)
	kmgTest.Equal(err,ErrStringWrongType)

	err=IncrBy("test_1",5)
	kmgTest.Equal(err,nil)
	kmgTest.Equal(MustGetIntIgnoreNotExist("test_1"),7)
}

func TestIncrByFloat(ot *testing.T){
	MustFlushDbV2()
	MustInsert("test_2","abc")
	MustZAdd("test_3",0,"abc")

	err:=IncrByFloat("test_1",2.1)
	kmgTest.Equal(err,nil)
	kmgTest.Equal(MustGetFloatIgnoreNotExist("test_1"),2.1)

	err = IncrByFloat("test_2",3)
	kmgTest.Equal(err,ErrValueNotFloatFormat)

	err = IncrByFloat("test_3",4)
	kmgTest.Equal(err,ErrStringWrongType)

	err=IncrByFloat("test_1",5.1)
	kmgTest.Equal(err,nil)
	kmgTest.Equal(MustGetFloatIgnoreNotExist("test_1"),7.2)

	err=IncrByFloat("test_1",1.7e200)
	kmgTest.Equal(err,nil)

	kmgTest.Equal(MustGetFloatIgnoreNotExist("test_1"),1.7e200)
}

func TestScanCallback(ot *testing.T){
	MustFlushDbV2()
	MustInsert("test_1","abc")
	MustInsert("test_2","abc")
	MustInsert("testno_3","abc")

	outKey:=[]string{}
	err := ScanCallback("*",func(key string) error{
		outKey = append(outKey,key)
		return nil
	})
	kmgTest.Equal(err,nil)
	kmgTest.Equal(len(outKey),3)
	kmgTest.Ok(kmgStrings.IsInSlice(outKey,"test_1"))
	kmgTest.Ok(kmgStrings.IsInSlice(outKey,"test_2"))
	kmgTest.Ok(kmgStrings.IsInSlice(outKey,"testno_3"))


	outKey=[]string{}
	err = ScanCallback("test_*",func(key string)error{
		outKey = append(outKey,key)
		return nil
	})
	kmgTest.Equal(err,nil)
	kmgTest.Equal(len(outKey),2)
	kmgTest.Ok(kmgStrings.IsInSlice(outKey,"test_1"))
	kmgTest.Ok(kmgStrings.IsInSlice(outKey,"test_2"))

	sList,err:=ScanWithOutputLimit("test_*",1)
	kmgTest.Equal(err,nil)
	kmgTest.Equal(len(sList),1)

	//benchmarkScanCallback()
}

func benchmarkScanCallback() {
	N:=1000*100
	MustFlushDbV2()
	t:=time.Now()
	pairList :=make([]KeyValuePair,N)
	for i:=0;i<N;i++ {
		pairList[i].Key = "test_"+strconv.Itoa(i)
		pairList[i].Value = "abc"
	}
	MustMSet(pairList)
	fmt.Println(time.Since(t)) //386.265848ms
	t = time.Now()
	num:=0
	err := ScanCallback("*",func(key string) error{
		num++
		return nil
	}) //169.983354ms
	fmt.Println(time.Since(t))
	kmgTest.Equal(err,nil)
	kmgTest.Equal(num,N)
	t = time.Now()
	sList,err:=Keys("*") //138.565292ms
	fmt.Println(time.Since(t))
	kmgTest.Equal(err,nil)
	kmgTest.Equal(len(sList),N)
	MustFlushDbV2()
}

/*
暂时不能使用这个办法来调试速度,容易爆内存
func BenchmarkScanCallback(b *testing.B) {
	b.StopTimer()
	MustFlushDbV2()
	fmt.Println(b.N)
	b.ResetTimer()
	pairList :=make([]KeyValuePair,b.N)
	for i:=0;i<b.N;i++ {
		pairList[i].Key = "test_"+strconv.Itoa(i)
		pairList[i].Value = "abc"
	}
	MustMSet(pairList)
	num:=0
	err := ScanCallback("*",func(key string){
		num++
	})
	kmgTest.Equal(err,nil)
}
*/