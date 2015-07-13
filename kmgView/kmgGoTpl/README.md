kmg golang template engine
=================
A PHP like template engine write for golang.

一个像PHP的golang模板引擎.

当前版本不稳定,随时会产生不向前兼容.

* 大幅度减少学习成本,你只需要知道,在 <? ?> 里面写golang语句, 在 <?= ?> 里面写需要渲染的golang表达式即可.
* .gotpl 表示普通文本模板,此处不做任何自动转义 .gotplhtml 表示html模板,此处做3种xss的自动转义.
* 报错有文件名和行号.
* 对html的3种xss转义自动支持(kmgXss.H kmxXss.Urlv kmgXss.Jsonv)

### example
template file.
```html
<?
package example
type Input struct{
	Name     string
	Value    string
	ShowName string
	Comment  string
	Need     bool
	ReadOnly bool
	Id       string
}
func tplInputString(config Input)string{
?>
<div class="form-group has-feedback">
    <label class="col-sm-2 control-label"><?=config.ShowName?>
    <? if config.Need{ ?>
        <span style="color:red">*</span>
    <? } ?>

    <div class="col-sm-8">
        <input type="text" autocomplete="off" class="form-control"
               <? if config.ReadOnly{ ?>readonly<? } ?>
               name="<?=config.Name?>"
        value="<?=config.Value?>"/>
        <span style="font-size:12px;color:red">
            <? if config.Comment!=""{ ?>
                提示: <?=config.Comment?>
            <? } ?>
        </span>
    </div>
</div>
<? }
?>
```

use the function `kmgGoTpl.MustBuildTplInDir("src/github.com/bronze1man/kmg/kmgView/testFile")` to compile the template in your make file.


### TODO
* 允许在golang的字符串,注释,部分出现语法引擎使用过的关键字(<?= <? ?> func { } import ( ) 等)
* 自动golang类型分析(确定是否直接渲染字符串,慢?)

### reference
* 大部分点子来源于 https://github.com/sipin/gorazor
