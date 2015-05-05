package kmgHttp

import "net/url"

//sort url query to unique it
// 这个东西是做什么用的?
// @deprecated
func SortUrlQuery(vurl *url.URL) {
	query := vurl.Query()
	vurl.RawQuery = query.Encode()
}
