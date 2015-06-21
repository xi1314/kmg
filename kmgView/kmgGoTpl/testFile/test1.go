package example

import (
	"bytes"
	"github.com/bronze1man/kmg/kmgXss"
)

type Input struct {
	Name     string
	Value    string
	ShowName string
	Comment  string
	Need     bool
	ReadOnly bool
	Id       string
}

func tplInputString(config Input) string { //函数开头(需要想办法判断当前后面的内容是否处于某个函数里面)
	var _buf bytes.Buffer
	_buf.WriteString(`
<div class="form-group has-feedback">
    <label class="col-sm-2 control-label">`)
	_buf.WriteString(config.ShowName)
	_buf.WriteString(`
    `)
	if config.Need {
		_buf.WriteString(`
        <span style="color:red">*</span>
    `)
	}
	_buf.WriteString(`

    <div class="col-sm-8">
        <input type="text" autocomplete="off" class="form-control"
               `)
	if config.ReadOnly {
		_buf.WriteString(`readonly`)
	}
	_buf.WriteString(`
               name="`)
	_buf.WriteString(kmgXss.H(config.Name))
	_buf.WriteString(`"
        value="`)
	_buf.WriteString(kmgXss.H(config.Value))
	_buf.WriteString(`"/>
        <span style="font-size:12px;color:red">
            `)
	if config.Comment != "" {
		_buf.WriteString(`
                提示: `)
		_buf.WriteString(kmgXss.H(config.Comment))
		_buf.WriteString(`
            `)
	}
	_buf.WriteString(`
        </span>
    </div>
</div>
`)
	return _buf.Bytes()
} //需要根据 {} 匹配出来这个是函数结束,并且在函数结束的地方加入 return buf.String()
