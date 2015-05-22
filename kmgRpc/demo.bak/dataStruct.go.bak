package demo

import "time"

type AutoRegisterRequest struct {
	DeviceId string
}
type VpnType string

const (
	// need client version >= 1.1
	VpnTypeIPSECXuathAndPsk VpnType = "IPSECXuathAndPsk"
	VPNTypeIke2Psk          VpnType = "Ike2Psk"

	// need client version >= 1.1 old name is Ike2EspAndPsk
	VPNTypeIke2EapAndPsk VpnType = "Ike2EapAndPsk"
	VPNTypeIke2Eap       VpnType = "Ike2Esp"
)

type VpnConfig struct {
	Type       VpnType
	ServerAddr string
	Username   string
	Password   string
	Psk        string
}

//给前台用户看的信息
type FrontUserInfo struct {
	Id                  string
	Sk                  string
	Username            string
	VpnConfig           *VpnConfig
	ProductList         []Product
	DataTransferForever int64
	DataTransferMonth   int64
	MonthPayStartTime   time.Time //月度数据流量开始时间
	MonthPayEndTime     time.Time //月度数据流量结束时间
}

type Product struct {
	IapId          string
	Name           string
	DataTransfer   int64
	Type           ProductType
	PriceString    string //这个需要由客户端自行拉取,目前服务器实现起来太花时间了. 这个放在这里占个坑位
	OldPriceString string
}

type ProductType string

const (
	ProductTypeForever ProductType = "Forever"
	ProductTypeMonth   ProductType = "Month"
)
