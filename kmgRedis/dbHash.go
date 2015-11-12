package kmgRedis

func HGetAll(key string) (map[string]string,error){
	return gClient.HGetAllMap(key).Result()
}

func MustHGetAll(key string) (map[string]string){
	h,err:=HGetAll(key)
	if err!=nil{
		panic(err)
	}
	return h
}

func MustHSet(key1 string,key2 string,value string){
	err:=gClient.HSet(key1,key2,value).Err()
	if err!=nil{
		panic(err)
	}
}