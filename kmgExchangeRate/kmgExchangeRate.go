package kmgExchangeRate

import (
	"fmt"
	"github.com/bronze1man/kmg/encoding/kmgJson"
	"github.com/bronze1man/kmg/kmgCache"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"github.com/bronze1man/kmg/kmgStrconv"
	"github.com/bronze1man/kmg/kmgStrings"
	"github.com/bronze1man/kmg/kmgTime"
	"time"
)

type yqlRateResponse struct {
	Query struct {
		Results struct {
			Rate struct {
				Rate string `json:"rate"`
			} `json:"rate"`
		} `json:"results"`
	} `json:query`
}

//缓存时间为1天
func MustGetExchangeRateWithCache(from string, to string) float64 {
	b, err := kmgCache.FileTtlCache("MustGetExchangeRateWithCache_"+from+"_"+to, func() (b []byte, ttl time.Duration, err error) {
		f := MustGetExchangeRate(from, to)
		return []byte(kmgStrconv.FormatFloat(f)), kmgTime.Day, nil
	})
	if err != nil {
		panic(err)
	}
	f, err := kmgStrconv.ParseFloat64(string(b))
	if err != nil {
		panic(err)
	}
	return f
}

func MustGetExchangeRate(from string, to string) float64 {
	if !kmgStrings.IsAllAphphabet(from) {
		panic(fmt.Errorf("[MustGetExchangeRate] fromName [%s] should like USD", from))
	}
	if !kmgStrings.IsAllAphphabet(to) {
		panic(fmt.Errorf("[MustGetExchangeRate] toName [%s] should like USD", to))
	}
	if from == to {
		panic(fmt.Errorf("[MustGetExchangeRate] fromName [%s] =toName [%s] "))
	}
	out := kmgHttp.MustUrlGetContent(kmgHttp.MustSetParameterMapToUrl("https://query.yahooapis.com/v1/public/yql", map[string]string{
		"q":        fmt.Sprintf(`select * from yahoo.finance.xchange where pair="%s%s"`, from, to),
		"format":   "json",
		"env":      "store://datatables.org/alltableswithkeys",
		"callback": "",
	}))
	//out should look like {"query":{"count":1,"created":"2015-06-01T02:37:44Z","lang":"en","results":{"rate":{"id":"JPYCNY","Name":"JPY/CNY","Rate":"0.0501","Date":"6/1/2015","Time":"3:37am","Ask":"0.0501","Bid":"0.0501"}}}}
	resp := yqlRateResponse{}
	kmgJson.MustUnmarshal(out, &resp)
	if resp.Query.Results.Rate.Rate == "N/A" {
		panic(fmt.Errorf("[MustGetExchangeRate] Currency [%s][%s] not exist", from, to)) //有一种货币不存在.
	}
	rate, err := kmgStrconv.ParseFloat64(resp.Query.Results.Rate.Rate)
	if err != nil {
		panic(err)
	}
	return rate
}
