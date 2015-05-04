package kmgHttp

import (
	"net/url"
	"strings"
)

// @deprecated
func AddParameterToUrl(urls string, key string, value string) (urlout string, err error) {
	u, err := url.Parse(urls)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Add(key, value)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func SetParameterToUrl(urls string, key string, value string) (urlout string, err error) {
	u, err := url.Parse(urls)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set(key, value)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func MustSetParameterToUrl(urls string, key string, value string) (urlout string) {
	u, err := SetParameterToUrl(urls, key, value)
	if err != nil {
		panic(err)
	}
	return u
}

func SetParameterMapToUrl(urls string, row map[string]string) (urlout string, err error) {
	u, err := url.Parse(urls)
	if err != nil {
		return
	}
	q := u.Query()
	for key, value := range row {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func MustSetParameterMapToUrl(urls string, row map[string]string) (urlout string) {
	u, err := SetParameterMapToUrl(urls, row)
	if err != nil {
		panic(err)
	}
	return u
}

func GetDomainName(url string) (domainName string, protocol string) {
	protocolList := []string{
		"http://",
		"https://",
		"file://",
		"ftp://",
	}
	for _, p := range protocolList {
		if strings.HasPrefix(url, p) {
			protocol = p
			break
		}
	}
	if protocol == "" {
		return
	}
	domainName = strings.Replace(url, protocol, "", -1)
	list := strings.Split(domainName, "/")
	domainName = list[0]
	return
}
