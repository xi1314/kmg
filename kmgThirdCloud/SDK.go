package kmgThirdCloud

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

type Instance struct {
	Ip   string
	Id   string //第三方只认他们自己产生的ID
	Name string
}

var _ SDK = (*RackspaceSDK)(nil)
var _ SDK = (*AliyunSDK)(nil)

//以下函数均需做成「同步」的形式
type SDK interface {
	//Region       string //实例所在地区
	//InstanceName string
	//ImageName    string //操作系统
	//FlavorName   string //实例配置，如 4 CPU / 4GB RAM
	CreateInstance() (ip string)
	DeleteInstance(ip string)
	RenameInstanceByIp(name, ip string)
	ListAllInstance() (ipInstanceMap map[string]Instance)
}

type SDKCache struct {
	cache SDK
	Init  func() SDK
}

func (sdkCache SDKCache) Get() SDK {
	if sdkCache.cache == nil {
		sdkCache.cache = sdkCache.Init()
	}
	return sdkCache.cache
}
