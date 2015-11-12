package kmgRedis

func SAdd(key string,value string) (err error){
	return gClient.SAdd(key,value).Err()
}

func MustSAdd(key string,value string) {
	err := SAdd(key,value)
	if err!=nil{
		panic(err)
	}
	return
}

func MustSMembers(key string) (sList []string){
	sList,err:=gClient.SMembers(key).Result()
	if err!=nil{
		panic(err)
	}
	return sList
}