package kmgHttp

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func NewHttpsCertNotCheckClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

var httpsNotCheckClient *http.Client
var httpsNotCheckClientOnce sync.Once

func GetHttpsCertNotCheckClient() *http.Client {
	httpsNotCheckClientOnce.Do(func() {
		httpsNotCheckClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	})
	return httpsNotCheckClient
}

//一个不验证证书,并且使用http代理的http客户端
// httpProxy look like http://127.0.0.1:9876
func MustNewTestClientWithHttpProxy(httpProxy string) *http.Client {
	var Proxy func(*http.Request) (*url.URL, error)
	if httpProxy != "" {
		u, err := url.Parse(httpProxy)
		if err != nil {
			panic(err)
		}
		Proxy = http.ProxyURL(u)
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Proxy: Proxy,
		},
	}
}

func NewTimeoutHttpClient(dur time.Duration) *http.Client {
	return &http.Client{
		Timeout: dur,
	}
}

func NewTimeoutNoKeepAliveHttpClient(dur time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
		Timeout: dur,
	}
}
