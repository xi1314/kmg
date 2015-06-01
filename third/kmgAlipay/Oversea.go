package kmgAlipay

import (
	"bytes"
	"fmt"
	"github.com/bronze1man/kmg/kmgCrypto"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"github.com/bronze1man/kmg/kmgStrconv"
	"github.com/bronze1man/kmg/kmgTime"
	"sort"
	"strings"
	"time"
)

//境外收单接口系列
// https://download.alipay.com/ui/doc/global/cross-border_website_payment.zip
// 看英文版,中文版是旧的,少一些功能
type OverseaTrade struct {
	PartnerId    string
	NotifyUrl    string
	ReturnUrl    string
	SecurityCode string
}

type OverseaTradePayRequest struct {
	Subject    string
	Body       string //可选
	OutTradeNo string
	Currency   string
	//TotalFee 和 RmbFee 二选一,必须选一个
	TotalFee            float64 //可选
	RmbFee              float64 //可选
	Supplier            string  //可选
	TimeoutRule         string  //可选
	SpecifiedPayChannel string  //可选
	SellerId            string  //可选
	SellerName          string  //可选
	SellerIndustry      string  //可选
}

//用户发起支付
func (ot *OverseaTrade) Pay(req *OverseaTradePayRequest) (url string) {
	query := map[string]string{
		"service":               "create_forex_trade",
		"parter":                ot.PartnerId,
		"notify_url":            ot.NotifyUrl,
		"return_url":            ot.ReturnUrl,
		"subject":               req.Subject,
		"body":                  req.Body,
		"out_trade_no":          req.OutTradeNo,
		"currency":              req.Currency,
		"supplier":              req.Supplier,
		"timeout_rule":          req.TimeoutRule,
		"specified_pay_channel": req.SpecifiedPayChannel,
		"seller_id":             req.SellerId,
		"seller_name":           req.SellerName,
		"seller_industry":       req.SellerIndustry,
	}
	if req.TotalFee != 0 {
		query["total_fee"] = kmgStrconv.FormatFloat(req.TotalFee)
	}
	if req.RmbFee != 0 {
		query["rmb_fee"] = kmgStrconv.FormatFloat(req.RmbFee)
	}
	ot.md5Sign(query)
	return kmgHttp.MustSetParameterMapToUrl("https://mapi.alipay.com/gateway.do", query)
}

type OverseaTradeStatus string

const (
	OverseaTradeStatusFinish OverseaTradeStatus = "TRADE_FINISHED"
	OverseaTradeStatusClose  OverseaTradeStatus = "TRADE_CLOSED"
)

type OverseaTradeReturnInfo struct {
	OutTradeNo  string
	Currency    string
	TotalFee    float64
	TradeStatus OverseaTradeStatus //一定是 OverseaTradeStatusFinish
	TradeNo     string
}

// 同步回调
func (ot *OverseaTrade) MustReturnPage(ctx *kmgHttp.Context) (info OverseaTradeReturnInfo) {
	var err error
	info.OutTradeNo = ctx.MustInStr("out_trade_no")
	info.Currency = ctx.MustInStr("currency")
	info.TotalFee, err = kmgStrconv.ParseFloat64(ctx.MustInStr("total_fee"))
	if err != nil {
		panic(err)
	}
	info.TradeStatus = OverseaTradeStatus(ctx.MustInStr("trade_status"))
	info.TradeNo = ctx.MustInStr("trade_no")
	return info
}

type OverseaTradeNotifyInfo struct {
	NotifyId    string
	NotifyTime  time.Time
	OutTradeNo  string
	TradeStatus OverseaTradeStatus
	TradeNo     string
	Currency    string
	TotalFee    float64
}

// 异步回调,请不要在f中输出任何支付串.
func (ot *OverseaTrade) MustNotifyAction(ctx *kmgHttp.Context, f func(info OverseaTradeNotifyInfo)) {
	var err error
	ctx.MustPost()
	info := OverseaTradeNotifyInfo{}
	info.NotifyId = ctx.MustInStr("notify_id")
	info.NotifyTime = kmgTime.MustFromMysqlFormatInLocation(ctx.MustInStr("notify_time"), kmgTime.BeijingZone)
	info.OutTradeNo = ctx.MustInStr("out_trade_no")

	info.Currency = ctx.MustInStr("currency")
	info.TotalFee, err = kmgStrconv.ParseFloat64(ctx.MustInStr("total_fee"))
	if err != nil {
		panic(err)
	}
	info.TradeStatus = OverseaTradeStatus(ctx.MustInStr("trade_status"))
	info.TradeNo = ctx.MustInStr("trade_no")
	err = ot.md5Verify(ctx.GetInMap())
	if err != nil {
		panic(err)
	}
	err = ot.VerifyNotify(info.NotifyId)
	if err != nil {
		panic(err)
	}
	f()
	ctx.WriteString("success")
}

// 通知验证接口
// 和支付宝手机接口一模一样.
func (ot *OverseaTrade) VerifyNotify(NotifyId string) (err error) {
	u := kmgHttp.MustSetParameterMapToUrl("https://mapi.alipay.com/gateway.do", map[string]string{
		"service":   "notify_verify",
		"partner":   ot.PartnerId,
		"notify_id": NotifyId,
	})
	content, err := kmgHttp.UrlGetContent(u)
	if err != nil {
		return nil
	}
	if !bytes.Equal(content, []byte(`true`)) {
		return fmt.Errorf("notify_id verify fail")
	}
	return nil
}

// TODO 批量退款 上传退款文件接口
// TODO 下载对账文件接口
// TODO 单笔退款接口
// TODO 退款撤销接口
// TODO 下载清算文件接口
// TODO 下载汇率文件接口
// TODO 单条交易查询接口
// TODO 会员共享 ID 接口

type kv struct {
	K string
	V string
}
type kvSorter []kv

func (l kvSorter) Len() int      { return len(l) }
func (l kvSorter) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l kvSorter) Less(i, j int) bool {
	if l[i].K == l[i].K {
		return l[i].V < l[i].V
	}
	return l[i].K < l[i].K
}

//这个函数使用md5方式对query添加签名,没有数据的参数会被删除,并且直接更新输入的query数组.
func (ot *OverseaTrade) md5Sign(query map[string]string) {
	kvList := make([]kv, len(query))
	i := 0
	for k, v := range query {
		if v == "" {
			continue
		}
		kvList[i] = kv{
			K: k,
			V: v,
		}
		i++
	}
	sort.Sort(kvSorter(kvList))

	toEncodeList := make([]string, len(kvList))
	for i, data := range kvList {
		toEncodeList[i] = data.K + `="` + data.V + `"`
	}
	toSign := strings.Join(toEncodeList, "&")
	signed := kmgCrypto.Md5Hex(toSign + ot.SecurityCode)
	query["sign_type"] = "MD5"
	query["sign"] = signed
	return
}

func (ot *OverseaTrade) md5Verify(query map[string]string) (err error) {
	signed := query["sign"]
	delete(query, "sign")
	ot.md5Sign(query)
	if signed != query["sign"] {
		return fmt.Errorf("[md5Verify] fail")
	}
	return nil
}
