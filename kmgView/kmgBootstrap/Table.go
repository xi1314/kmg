package kmgBootstrap

func NewTableRender(wrapName string,tableName string,table Table)string{
	return NewWrap(
		wrapName,
		Panel{
			Title:  tableName,
			Body:   table,
		},
	).HtmlRender()
}
