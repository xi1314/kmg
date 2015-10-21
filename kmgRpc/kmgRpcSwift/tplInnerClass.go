package kmgRpcSwift

import (
	"bytes"
)

func (config InnerClass) tplInnerClass() string {
	var _buf bytes.Buffer
	_buf.WriteString(`
    `)
	if config.IsPublic {
	} else {
		_buf.WriteString(`private`)
	}
	_buf.WriteString(` struct `)
	_buf.WriteString(config.Name)
	_buf.WriteString(`{
        `)
	for _, field := range config.FieldList {
		_buf.WriteString(`
            var `)
		_buf.WriteString(field.Name)
		_buf.WriteString(`:`)
		_buf.WriteString(field.TypeStr)
		_buf.WriteString(` = `)
		_buf.WriteString(field.TypeStr)
		_buf.WriteString(`()
        `)
	}
	_buf.WriteString(`
`)
	if config.IsPublic {
		_buf.WriteString(`
mutating func ToData(inData:JSON){
`)
		for _, field := range config.FieldList {
			_buf.WriteString(`
`)
			if field.TypeStr == "Int" || field.TypeStr == "NSString" {
				_buf.WriteString(`
self.`)
				_buf.WriteString(field.Name)
				_buf.WriteString(` = inData["Out_0"]["`)
				_buf.WriteString(field.Name)
				_buf.WriteString(`"].`)
				switch field.TypeStr {
				case "Int":
					_buf.WriteString(`intValue`)
				case "NSString":
					_buf.WriteString(`stringValue`)
				default:
					_buf.WriteString(`stringValue`)
				}
				_buf.WriteString(`
`)
			} else {
				_buf.WriteString(`
self.`)
				_buf.WriteString(field.Name)
				_buf.WriteString(`.ToData(inData)
`)
			}
			_buf.WriteString(`
`)
		}
		_buf.WriteString(`
}
`)
	}
	_buf.WriteString(`
    }
`)
	return _buf.String()
}
