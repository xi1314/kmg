kmgViewResource web前端资源管理工具
=====================================

### js和css,图片,swf等资源集成方法
* 在kmgViewResource 中资源被划分为不同的package,比如 bootstrap 就被划分为 名叫 github.com/bronze1man/kmg/kmgView/kmgWeb/bootstrap 的package.
* 每一个package占用一个目录
    * 这个目录里面所有的文件都是这个package的一部分.但是这个目录的子目录里面的文件不算这个package的一部分.
    * 比如 名为 github.com/bronze1man/kmg/kmgView/kmgWeb/bootstrap 的package 就会被放置到 src/github.com/bronze1man/kmg/kmgView/kmgWeb/bootstrap 这个位置.
* package里面可以放置 js文件
    * js文件会按照package里面的文件的名称进行排序,然后按照此顺序合并到一个js文件里面.
* package里面可以放置 css文件
    * css文件会按照package里面的文件的名称进行排序,然后按照此顺序合并到一个css文件里面.
* package里面可以放置 图片,字体和swf文件
    * 图片和swf会生成到网站的一个资源根目录下
    * 图片和swf的文件路径会本忽略掉,只保留文件名称.
    * 使用 getBootstrapViewResource().GetUrlPrefix() 获取本次打包的资源根目录.
    * 比如:
        * 资源文件 src/github.com/bronze1man/kmg/kmgView/kmgWeb/font-awesome/FontAwesome.otf 会被生成到
        getBootstrapViewResource().GetUrlPrefix()+"/"+FontAwesome.otf 这个url位置.
        * 在js里面使用 getResourceUrlPrefix()+"/"+FontAwesome.otf 引用到这个文件
        * 在css里面使用 FontAwesome.otf 引用这个文件 (css和字体都在资源根目录下)
        * 在golang里面使用 getBootstrapViewResource().GetUrlPrefix()+"/"+FontAwesome.otf 引用到这个文件

* 每一个package依赖其他package.
    * 在 package 里面 任意js和css前面加上类似与下面的注释:
```
/*
import(
    "github.com/bronze1man/kmg/kmgView/kmgWeb/jquery"
    "github.com/bronze1man/kmg/kmgView/kmgWeb/bootstrap"
)
 */
```
    * 被依赖的pacakge的js和css在合并之后的js和css一定在依赖的package之前.
    * 不允许循环依赖.
    * 比如 bootstrap 依赖 jquery
        * 在 src/github.com/bronze1man/kmg/kmgView/kmgWeb/bootstrap 目录下面建立一个 import.js 文件,里面放入如下内容:
```
/*
import(
    "github.com/bronze1man/kmg/kmgView/kmgWeb/jquery"
)
 */
```
        * 编译器保证jquery 的js一定在bootstrap的js之前.


### build部分使用方法
* 写一个更新上传资源的脚本,例如:
```
	kmgViewResource.ResourceBuild(&kmgViewResource.ResourceUploadRequest{
		ImportPathList: []string{
			"github.com/bronze1man/kmg/kmgView/kmgWeb/bootstrap",
			"github.com/bronze1man/kmg/kmgView/kmgWeb/font-awesome",
			"github.com/bronze1man/kmg/kmgView/kmgWeb/jquery",
			"github.com/bronze1man/kmg/kmgView/kmgWeb/moment",
		},
		Qiniu:         getKmgToolQiniu(),
		QiniuPrefix:   "kmgBootstrap",
		OutGoFilePath: "src/github.com/bronze1man/kmg/kmgView/kmgBootstrap/generated_BootstrapResource.go",
		Name:    "Bootstrap",
	})
```

* 在html模板的开头插入 header代码, 在结束位置插入 footer代码
```
<? package kmgBootstrap
func tplWrap (w Wrap) string { ?>
<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title><?= w.Title ?></title>
    <?= raw(getBootstrapViewResource().HeaderHtml()) ?>
</head>
<body style="padding: 20px;">
    <?= raw(getBootstrapViewResource().FooterHtml()) ?>
</body>
</html>
<? } ?>
```
* 如果传入qiniu相关的信息,生成代码时,会向七牛上传资源文件.此时可以生成二进制文件,然后在没有项目源代码的情况使用本资源管理器.
* 如果不传入qiniu相关的信息,生成后的资源只能在有项目源代码的情况下使用. (TODO finish it)
* 在 github.com/bronze1man/kmg/kmgView/kmgBootstrap 这个package 里面使用 getBootstrapViewResource() 在golang里面访问刚才那个例子里面生成的前端资源管理对象.
* 在第一次访问前端资源管理对象时,进行初始化.
* 有项目源代码时(以存在.kmg.yml为准),此时 前端资源管理对象 处于调试模式.每次访问HeaderHtml会重新编译一次.
* 没有项目源代码时,此时 前端资源管理对象 处于上线模式.使用一个到七牛的反向代理提供前端文件.

### 作用
* 使用类似golang的依赖树的方式,对前端的 js,css,图片,swf 等进行管理.
* 大幅度降低前端资源管理时的,配置复杂度.
* 支持开发时,文件变化自动检测,减少开发时的复杂度.(或者说需要考虑的东西)