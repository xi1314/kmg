package kmgLog

// 如果f发生错误,写一条log
// @deprecated
func LogErrCallback(category string, context interface{}, f func() error) error {
	err := f()
	if err == nil {
		return nil
	}
	Log(category, err.Error(), context)
	return err
}
