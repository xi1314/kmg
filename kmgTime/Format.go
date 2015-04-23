package kmgTime

import "time"

//输出成mysql的格式,并且使用默认时区,并且在0值的时候输出空字符串
func DefaultFormat(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(DefaultTimeZone).Format(FormatMysql)
}
