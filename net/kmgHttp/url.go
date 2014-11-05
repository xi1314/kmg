package kmgHttp

import (
	"net/url"
)

func AddParameterToUrl(urls string, s map[string]string) (urlout string, err error) {
	u, err := url.Parse(urls)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for k, v := range s {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}
