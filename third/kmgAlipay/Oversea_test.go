package kmgAlipay

import (
	"github.com/bronze1man/kmg/kmgTest"
	"strings"
	"testing"
)

func TestOverseaMd5Sign(ot *testing.T) {
	trade := OverseaTrade{
		SecurityCode: "56Tae2ZROl2DSw",
	}
	query := map[string]string{
		"service":        "create_forex_trade",
		"partner":        "20881234567890123",
		"return_url":     "http://xxxx.com/?n=Xxx.Pay.ReturnPage",
		"notify_url":     "http://xxxx.com/?n=Xxx.Pay.NotifyAction",
		"_input_charset": "utf-8",
		"subject":        "多项测试",
		"body":           "test哈哈only",
		"out_trade_no":   "1433229365",
		"total_fee":      "10",
		"currency":       "JPY",
	}
	trade.md5Sign(query)
	kmgTest.Equal(query["sign"], "f07ac88d67becf081a06baa3a13656a8")
}

func TestOverseaPay(ot *testing.T) {
	trade := OverseaTrade{
		PartnerId:    "20881234567890123",
		SecurityCode: "56Tae2ZROl2DSw",
	}
	url := trade.Pay(&OverseaTradePayRequest{
		Subject:    "多项测试",
		Body:       "test哈哈only",
		OutTradeNo: "1433229365",
		Currency:   "JPY",
		TotalFee:   10,
	})
	kmgTest.Ok(strings.Contains(url, "8d927a97f15a2f40a31040267a1d9afe"))
	//	kmgTest.Equal(query["sign"],"b49d7f5e6341e66870473222edc5df0b")
}

func TestOverseaMd5Verify(ot *testing.T) {
	trade := OverseaTrade{
		SecurityCode: "56Tae2ZROl2DSw",
	}
	query := map[string]string{
		"service":        "create_forex_trade",
		"partner":        "20881234567890123",
		"return_url":     "http://xxxx.com/?n=Xxx.Pay.ReturnPage",
		"notify_url":     "http://xxxx.com/?n=Xxx.Pay.NotifyAction",
		"_input_charset": "utf-8",
		"subject":        "多项测试",
		"body":           "test哈哈only",
		"out_trade_no":   "1433229365",
		"total_fee":      "10",
		"currency":       "JPY",
		"sign":           "f07ac88d67becf081a06baa3a13656a8",
		"sign_type":      "MD5",
	}
	err := trade.md5Verify(query)
	kmgTest.Equal(err, nil)

	trade = OverseaTrade{
		SecurityCode: "56Tae2ZROl2DSw",
	}
	query = map[string]string{
		"currency": "JPY",
		//"n": "Sig.Front.Pay.ReturnPage",
		"out_trade_no": "1433229365",
		"sign":         "bab35fb00e9e858a0265241a94375ec2",
		"sign_type":    "MD5",
		"total_fee":    "1.00",
		"trade_no":     "2015060200001000110056045982",
		"trade_status": "TRADE_FINISHED",
	}
	err = trade.md5Verify(query)
	kmgTest.Equal(err, nil)
}

//func TestOverseaMustReturnPage(t *testing.T) {
//	trade := OverseaTrade{
//		PartnerId: "20881234567890123",
//		SecurityCode: "56Tae2ZROl2DSw",
//	}
//	ctx := kmgHttp.NewTestContext().
//	SetInStr("out_trade_no", "1433229365").
//	SetInStr("currency", "JPY").
//	SetInStr("total_fee","1.00").
//	SetInStr("sign","bab35fb00e9e858a0265241a94375ec2").
//	SetInStr("sign_type","MD5").
//	SetInStr("trade_status", "TRADE_FINISHED").
//	SetInStr("trade_no", "2015060200001000110056045982").
//	SetPost()
//	info := trade.MustReturnPage(ctx)
//	kmgTest.Ok(strings.Contains(info.TradeNo, "1433229365"))
//}
