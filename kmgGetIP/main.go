package kmgGetIP

import (
	"github.com/bronze1man/kmg/kmgErr"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"io/ioutil"
	"net/http"
	"strings"
)

//Server Side
func GetIPFromCTXAndRespond(ctx *kmgHttp.Context) {
	req := ctx.GetRequest()
	address := req.RemoteAddr
	address = strings.Split(address, ":")[0]
	ctx.WriteString("OK" + address)
}

//Client Side
func GetIPFromRemote(remoteServerUrl string) (ip string) {
	resp, err := http.Get(remoteServerUrl)
	kmgErr.PanicIfError(err)
	body, err := ioutil.ReadAll(resp.Body)
	kmgErr.PanicIfError(err)
	ip = string(body)
	if !strings.HasPrefix(ip, "OK") {
		panic("WhoAmI Error: can not get a Ipv4 address " + ip)
	}
	return strings.TrimPrefix(ip, "OK")
}
