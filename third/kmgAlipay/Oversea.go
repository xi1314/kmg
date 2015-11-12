package kmgAlipay

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bronze1man/kmg/kmgControllerRunner"
	"github.com/bronze1man/kmg/kmgCrypto"
	"github.com/bronze1man/kmg/kmgLog"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"github.com/bronze1man/kmg/kmgStrconv"
)

/*
境外收单接口系列
 文档在此 https://download.alipay.com/ui/doc/global/cross-border_website_payment.zip
 看英文版,中文版是旧的,少一些功能
 已经出现的错误,及其可能原因:
 * 不能重复创建交易，请返回订单页面重新付款，或重新登录支付宝付款。 同一个订单号,如果已经被支付成功了,不能再被创建支付了

*/

// 请调用Init初始化一下
// 请在进程初始化的时候进行Init.不要做懒加载.
type OverseaTrade struct {
	PartnerId    string
	SecurityCode string

	//自己网站的scheme和host,在支付回调处使用. 例如 https://www.abc.com ,注意最后不要加入 /
	SelfSchemeAndHost string
	// 支付处理回调, 支付成功必须处理,关闭可选处理.
	// @deprecated 请使用  PayFinishCallback 和 PayCloseCallback ,这个接口容易丢掉 是成功还是关闭的处理.
	PayCallback func(info OverseaTradeTransaction) (err error)
	// 支付成功处理回调. 出现错误,请直接panic
	PayFinishCallback func(info OverseaTradeTransaction)
	// 交易关闭处理回调. 出现问题,请直接panic
	PayCloseCallback func(info OverseaTradeTransaction)
	//支付返回页面处理回调,此时已经调用过支付处理回调了.
	PayReturnPageCallback func(info OverseaTradeTransaction, ctx *kmgHttp.Context)
}

// @deprecated 请使用 InitForPayCallback 或不初始化.
func (ot *OverseaTrade) Init() {
	ot.InitForPayCallback()
}

// 请在进程初始化的时候进行Init.不要做懒加载.避免掉单.
func (ot *OverseaTrade) InitForPayCallback() {
	if ot.SelfSchemeAndHost == "" {
		panic("支付回调必须填写当前网站scheme和host 如: http://127.0.0.1")
	}
	if ot.PayReturnPageCallback == nil {
		panic("支付回调必须处理同步回调.")
	}
	if ot.PayFinishCallback == nil && ot.PayCallback == nil {
		panic("支付成功必须处理,必须加入支付成功回调.")
	}
	kmgControllerRunner.RegisterController(ot)
}

// 注意事项:
// 1.不管是使用RmbFee 还是 TotalFee 如果Currency的值是日元,日元金额不能少于1JPY.
// 2.同一个订单号,如果已经被支付成功了,不能再被创建支付了
// 3.同一个订单号,如果没有支付过,可以随意修改价钱,然后可以再重新提交.
type OverseaTradePayRequest struct {
	Subject    string
	Body       string //可选
	OutTradeNo string // 最长64个字节
	Currency   string
	//TotalFee 和 RmbFee 二选一,必须选一个
	TotalFee            float64
	RmbFee              float64 //选rmb也需要传入Currency参数,Currency依然填对应的货币,支付宝会反向再算出那个货币需要的金额.
	Supplier            string  //可选
	TimeoutRule         string  //可选
	SpecifiedPayChannel string  //可选
	SellerId            string  //可选
	SellerName          string  //可选
	SellerIndustry      string  //可选

	// (低级接口)回调url,可以添加query参数.在回调里面请先删除这些query参数, 这个会覆盖掉OverseaTrade里面的那个使用kmgControllerRunner的配置
	NotifyUrl string
	ReturnUrl string
}

//用户发起支付,请传入Request,然后redirect到这个函数返回的url里面去.
func (ot *OverseaTrade) Pay(req *OverseaTradePayRequest) (url string) {
	query := map[string]string{
		"_input_charset":        "utf-8",
		"service":               "create_forex_trade",
		"partner":               ot.PartnerId,
		"notify_url":            ot.SelfSchemeAndHost + "/?n=github.com.bronze1man.kmg.third.kmgAlipay.OverseaTrade.NotifyAction",
		"return_url":            ot.SelfSchemeAndHost + "/?n=github.com.bronze1man.kmg.third.kmgAlipay.OverseaTrade.ReturnPage",
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
	if req.NotifyUrl != "" {
		query["notify_url"] = req.NotifyUrl
	}
	if req.ReturnUrl != "" {
		query["return_url"] = req.ReturnUrl
	}
	if req.TotalFee != 0 {
		if req.Currency == "JPY" {
			//日元精度是整数,传入小数点会报错,这个文档上没有写,但是经过实验发现是这个样子的.
			query["total_fee"] = kmgStrconv.FormatFloatPrec0(req.TotalFee)
		} else {
			query["total_fee"] = kmgStrconv.FormatFloatPrec2(req.TotalFee)
		}
	}
	if req.RmbFee != 0 {
		query["rmb_fee"] = kmgStrconv.FormatFloatPrec2(req.RmbFee)
	}
	ot.md5Sign(query)
	kmgLog.Log("Alipay", "Oversea Pay", query)
	//已经手动验证url里面的参数顺序无关紧要.
	return kmgHttp.MustSetParameterMapToUrl("https://mapi.alipay.com/gateway.do", query)
}

// 请不要手动调用,这个是自动注册到 kmgControllerRunner里面的
func (ot *OverseaTrade) ReturnPage(ctx *kmgHttp.Context) {
	ctx.DeleteInMap("n")
	info := ot.MustReturnPage(ctx)
	ot.payCallbackProceess(info)
	ot.PayReturnPageCallback(info, ctx)
}

// 请不要手动调用,这个是自动注册到 kmgControllerRunner里面的
func (ot *OverseaTrade) NotifyAction(ctx *kmgHttp.Context) {
	ctx.DeleteInMap("n")
	ot.mustNotifyActionV2(ctx, ot.payCallbackProceess)
}

func (ot *OverseaTrade) payCallbackProceess(info OverseaTradeTransaction) {
	if ot.PayCallback != nil {
		err := ot.PayCallback(info)
		if err != nil {
			panic(err)
		}
	}
	if info.TradeStatus == OverseaTradeStatusFinish && ot.PayFinishCallback != nil {
		ot.PayFinishCallback(info)
	}
	if info.TradeStatus == OverseaTradeStatusClose && ot.PayCloseCallback != nil {
		ot.PayCloseCallback(info)
	}
}

type OverseaTradeStatus string

const (
	OverseaTradeStatusFinish OverseaTradeStatus = "TRADE_FINISHED"
	OverseaTradeStatusClose  OverseaTradeStatus = "TRADE_CLOSED"
)

//一条交易信息
type OverseaTradeTransaction struct {
	OutTradeNo  string //用户传入的 OutTradeNo
	Currency    string
	TotalFee    float64            //此处是外币的价格,即使传入RmbFee,此处也是外币的价格
	TradeStatus OverseaTradeStatus //一定是 OverseaTradeStatusFinish
	TradeNo     string             //支付宝id
	Subject     string
}

// 同步回调
// 调用前请清除您自己的参数.
// @deprecated 请使用 OverseaTrade.PayFinishCallback 和 OverseaTrade.PayCloseCallback
func (ot *OverseaTrade) MustReturnPage(ctx *kmgHttp.Context) (info OverseaTradeTransaction) {
	kmgLog.Log("Alipay", "Oversea PayReturnPage", ctx.GetInMap())
	var err error
	info.OutTradeNo = ctx.MustInStr("out_trade_no")
	info.Currency = ctx.MustInStr("currency")
	info.TotalFee, err = kmgStrconv.ParseFloat64(ctx.MustInStr("total_fee"))
	if err != nil {
		panic(err)
	}
	info.TradeStatus = OverseaTradeStatus(ctx.MustInStr("trade_status"))
	info.TradeNo = ctx.MustInStr("trade_no")
	//这个也可以验证数据,只是文档上面没写.
	err = ot.md5Verify(ctx.GetInMap())
	if err != nil {
		panic(err)
	}
	// 向支付宝询问这个订单的情况
	oInfo := ot.MustSingleTransactionQuery(info.OutTradeNo)
	if oInfo.TradeStatus != info.TradeStatus {
		panic("两次查询订单状态不一致")
	}
	info.Subject = oInfo.Subject
	return info
}

// 异步回调,请不要在f中输出任何支付串.(字符串？)
// 调用前请清除您自己的参数.
// @deprecated 请使用 OverseaTrade.PayFinishCallback 和 OverseaTrade.PayCloseCallback
func (ot *OverseaTrade) MustNotifyAction(ctx *kmgHttp.Context, f func(info OverseaTradeTransaction) (err error)) {
	ot.mustNotifyActionV2(ctx, func(info OverseaTradeTransaction) {
		err := f(info)
		if err != nil {
			panic(err)
		}
	})
}

func (ot *OverseaTrade) mustNotifyActionV2(ctx *kmgHttp.Context, f func(info OverseaTradeTransaction)) {
	kmgLog.Log("Alipay", "Oversea PayNotifyAction", ctx.GetInMap())
	var err error
	ctx.MustPost()
	info := OverseaTradeTransaction{}
	//info.NotifyId = ctx.MustInStr("notify_id") 这两项没有什么意义.
	//info.NotifyTime = kmgTime.MustFromMysqlFormatInLocation(ctx.MustInStr("notify_time"), kmgTime.BeijingZone)
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
	err = ot.VerifyNotify(ctx.MustInStr("notify_id"))
	if err != nil {
		panic(err)
	}
	// 向支付宝询问这个订单的情况
	oInfo := ot.MustSingleTransactionQuery(info.OutTradeNo)
	if oInfo.TradeStatus != info.TradeStatus {
		panic("两次查询订单状态不一致")
	}
	info.Subject = oInfo.Subject
	f(info)
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
		return err
	}
	if !bytes.Equal(content, []byte(`true`)) {
		return fmt.Errorf("notify_id verify fail")
	}
	return nil
}

type ExchangeRate struct {
	Time     time.Time
	Currency string
	Rate     float64
}

//下载汇率文件接口
func (ot *OverseaTrade) MustGetExchangeRateList() (output []ExchangeRate) {
	query := map[string]string{
		"service": "forex_rate_file",
		"partner": ot.PartnerId,
	}
	ot.md5Sign(query)
	content := kmgHttp.MustUrlGetContent(kmgHttp.MustSetParameterMapToUrl("https://mapi.alipay.com/gateway.do", query))

	lineList := strings.Split(string(content), "\n")
	output = make([]ExchangeRate, 0, len(lineList))
	for _, line := range lineList {
		line := strings.TrimSpace(line)
		if line == "" {
			continue
		}
		part := strings.Split(line, "|")
		if len(part) < 4 {
			panic(fmt.Errorf("[MustGetExchangeRateList] format error"))
		}
		t, err := time.Parse("20060102150405", part[0]+part[1]) //考虑别处基本不会使用,就直接写在这个地方了.
		if err != nil {
			panic(err)
		}
		rate, err := kmgStrconv.ParseFloat64(part[3])
		if err != nil {
			panic(err)
		}
		output = append(output, ExchangeRate{
			Time:     t,
			Currency: part[2],
			Rate:     rate,
		})
	}
	return output
}

type overseaTradeTransactionQueryResponse struct {
	XMLName     xml.Name           `xml:"alipay"`
	IsSuccess   string             `xml:"is_success"`
	TradeNo     string             `xml:"response>trade>trade_no"`
	OutTradeNo  string             `xml:"response>trade>out_trade_no"`
	Subject     string             `xml:"response>trade>subject"`
	TradeStatus OverseaTradeStatus `xml:"response>trade>trade_status"`
	Error       string             `xml:"error"`
}

func (ot *OverseaTrade) MustSingleTransactionQuery(outTradeId string) *OverseaTradeTransaction {
	tran, err := ot.SingleTransactionQuery(outTradeId)
	if err != nil {
		panic(err)
	}
	return tran
}

// 给调用者mock用
type SingleTransactionQueryer func(outTradeId string) (tran *OverseaTradeTransaction, err error)

// 单条交易查询接口
func (ot *OverseaTrade) SingleTransactionQuery(outTradeId string) (tran *OverseaTradeTransaction, err error) {
	query := map[string]string{
		"service":        "single_trade_query",
		"partner":        ot.PartnerId,
		"_input_charset": "utf-8",
		"out_trade_no":   outTradeId,
	}
	ot.md5Sign(query)
	content := kmgHttp.MustUrlGetContent(kmgHttp.MustSetParameterMapToUrl("https://mapi.alipay.com/gateway.do", query))
	response := overseaTradeTransactionQueryResponse{}
	err = xml.Unmarshal(content, &response)
	if err != nil {
		return nil, err
	}
	if response.IsSuccess != "T" {
		return nil, fmt.Errorf("[支付宝单条交易查询接口错误] [%s]", response.Error)
	}
	return &OverseaTradeTransaction{
		TradeNo:     response.TradeNo,
		OutTradeNo:  response.OutTradeNo,
		Subject:     response.Subject,
		TradeStatus: response.TradeStatus,
	}, nil
}

// TODO 批量退款 上传退款文件接口
// TODO 下载对账文件接口
// TODO 单笔退款接口
// TODO 退款撤销接口
// TODO 下载清算文件接口
// TODO 会员共享 ID 接口

type kv struct {
	K string
	V string
}
type kvSorter []kv

func (l kvSorter) Len() int      { return len(l) }
func (l kvSorter) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l kvSorter) Less(i, j int) bool {
	if l[i].K == l[j].K {
		return l[i].V < l[j].V
	}
	ret := l[i].K < l[j].K
	return ret
}

//这个函数使用md5方式对query添加签名,没有数据的参数会被删除,并且直接更新输入的query数组.
// 和手机app版细节有所不同.
func (ot *OverseaTrade) md5Sign(query map[string]string) {
	kvList := make([]kv, 0, len(query))
	for k, v := range query {
		if v == "" {
			delete(query, k)
			continue
		}
		kvList = append(kvList, kv{
			K: k,
			V: v,
		})
	}
	sort.Sort(kvSorter(kvList))

	toEncodeList := make([]string, len(kvList))
	for i, data := range kvList {
		toEncodeList[i] = data.K + `=` + data.V
	}
	toSign := strings.Join(toEncodeList, "&")
	//fmt.Println(kmgBase64.Base64EncodeStringToString(toSign))
	signed := kmgCrypto.Md5Hex([]byte(toSign + ot.SecurityCode))
	query["sign_type"] = "MD5"
	query["sign"] = signed
	return
}

func (ot *OverseaTrade) md5Verify(query map[string]string) (err error) {
	signed := query["sign"]
	delete(query, "sign")
	delete(query, "sign_type")
	ot.md5Sign(query)
	if signed != query["sign"] {
		return fmt.Errorf("[md5Verify] fail")
	}
	return nil
}
