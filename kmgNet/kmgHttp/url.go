package kmgHttp

import (
	"net/url"
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
