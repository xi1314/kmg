kmg golang template engine
=================
A PHP like template engine write for golang.
一个像PHP的golang模板引擎.

### example
template file.
```
<?
package example
import (
    "github.com/bronze1man/kmg/kmgXss"
)
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
               name="<?=kmgXss.H(config.Name)?>"
        value="<?=kmgXss.H(config.Value)?>"/>
        <span style="font-size:12px;color:red">
            <? if config.Comment!=""{ ?>
                提示: <?=kmgXss.H(config.Comment)?>
            <? } ?>
        </span>
    </div>
</div>
<? }
?>
```
output


### TODO
* 对html的3种xss转义自动支持(kmgXss.H kmxXss.Urlv kmgXss.Jsonv)
* 自动golang类型分析(确定是否直接渲染字符串,慢?)
* 允许在golang的字符串,注释,部分出现语法引擎使用过的关键字(<?= <? ?> func { } import ( ) 等)

### reference
* the idea of this type template engine is come from https://github.com/sipin/gorazor