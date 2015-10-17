package kmgViewResource

import (
	"github.com/bronze1man/kmg/encoding/kmgJson"
	"github.com/bronze1man/kmg/kmgCache"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"github.com/bronze1man/kmg/kmgXss"
	"strings"
	"sync"
	"reflect"
)

type Generated struct {
	Name                 string //拿去做缓存控制的.
	GeneratedJsFileUrl   string
	GeneratedCssFileUrl  string
	GeneratedJsFileName  string
	GeneratedCssFileName string
	GeneratedUrlPrefix   string // 末尾不包含 /
	RequestImportList    []string

	locker     sync.Mutex
	cachedInfo resourceBuildToDirResponse
	initOnce sync.Once
}

type htmlTplData struct {
	urlPrefix  string
	jsFileUrl  string
	cssFileUrl string
}

func (g *Generated) HtmlRender() string {
	data := g.GetHtmlTplData()
	return `<script>
		function getResourceUrlPrefix(){
			return ` + kmgXss.Jsonv(data.urlPrefix) + `
		}
	</script>
	<link rel="stylesheet" href="` + kmgXss.H(data.cssFileUrl) + `">
	<script src="` + kmgXss.H(data.jsFileUrl) + `"></script>`
}

func (g *Generated) HeaderHtml() string {
	data := g.GetHtmlTplData()
	return `<script>
		function getResourceUrlPrefix(){
			return ` + kmgXss.Jsonv(data.urlPrefix) + `
		}
	</script>
	<link rel="stylesheet" href="` + kmgXss.H(data.cssFileUrl) + `">`
}

func (g *Generated) FooterHtml() string {
	data := g.GetHtmlTplData()
	return `<script src="` + kmgXss.H(data.jsFileUrl) + `"></script>`
}

func (g *Generated) GetHtmlTplData() htmlTplData {
	if kmgConfig.HasDefaultEnv() {
		g.recheckAndReloadCache()
		g.locker.Lock()
		cachedInfo := g.cachedInfo
		g.locker.Unlock()
		urlPrefix := "/kmgViewResource." + g.Name
		return htmlTplData{
			urlPrefix:  urlPrefix,
			jsFileUrl:  urlPrefix + "/" + cachedInfo.JsFileName,
			cssFileUrl: urlPrefix + "/" + cachedInfo.CssFileName,
		}
	} else {
		urlPrefix := "/kmgViewResource." + g.Name
		return htmlTplData{
			urlPrefix:  urlPrefix,
			jsFileUrl:  urlPrefix + "/" + g.GeneratedJsFileName,
			cssFileUrl: urlPrefix + "/" + g.GeneratedCssFileName,
		}
	}
}

// 返回url前缀,末尾不包含 '/'
func (g *Generated) GetUrlPrefix() string {
	if kmgConfig.HasDefaultEnv() {
		g.recheckAndReloadCache()
		return "/kmgViewResource." + g.Name
	} else {
		return "/kmgViewResource." + g.Name
	}
}

func (g *Generated) GetJsUrl()string{
	if kmgConfig.HasDefaultEnv() {
		g.recheckAndReloadCache()
		g.locker.Lock()
		cachedInfo := g.cachedInfo
		g.locker.Unlock()
		return "/kmgViewResource." + g.Name+"/"+cachedInfo.JsFileName
	} else {
		return "/kmgViewResource." + g.Name+ "/" + g.GeneratedJsFileName
	}
}
func (g *Generated) GetCssUrl()string{
	if kmgConfig.HasDefaultEnv() {
		g.recheckAndReloadCache()
		g.locker.Lock()
		cachedInfo := g.cachedInfo
		g.locker.Unlock()
		return "/kmgViewResource." + g.Name+"/"+cachedInfo.CssFileName
	} else {
		return "/kmgViewResource." + g.Name+ "/" + g.GeneratedCssFileName
	}
}

// 这个初始化会在第一次使用的时候,自动进行,如果嫌自动初始化太慢,可以手动初始化.
func (g *Generated) Init() {
	g.initOnce.Do(func(){
		if kmgConfig.HasDefaultEnv() {
			g.recheckAndReloadCache()
			kmgHttp.MustAddFileToHttpPathToDefaultServer("/kmgViewResource."+g.Name+"/",
				kmgConfig.DefaultEnv().PathInProject("tmp/kmgViewResource_debug/"+g.Name))
		} else {
			// 默认使用反向代理方式提供数据.
			kmgHttp.MustAddUriProxyRefToUriToDefaultServer("/kmgViewResource."+g.Name+"/", g.GeneratedUrlPrefix)
		}
	})
}

// 获取某个资源文件的内容
func (g *Generated) GetContentByName(name string) (b []byte, err error) {
	name = strings.TrimPrefix(name, "/")
	if kmgConfig.HasDefaultEnv() {
		g.recheckAndReloadCache()
		path := kmgConfig.DefaultEnv().PathInProject("tmp/kmgViewResource_debug/" + g.Name + "/" + name)
		return kmgFile.ReadFile(path)
	} else {
		// 默认使用反向代理方式提供数据.
		return kmgHttp.UrlGetContent(g.GeneratedUrlPrefix + "/" + name)
	}

}


func (g *Generated) recheckAndReloadCache() {
	// 加载缓存文件,确定有哪些文件需要检查.
	cachedInfo := &resourceBuildToDirResponse{}
	err := kmgJson.ReadFile(kmgConfig.DefaultEnv().PathInProject("tmp/kmgViewResource_meta/"+g.Name), &cachedInfo)
	if err != nil || cachedInfo.JsFileName == "" {
		g.reloadCache()
		return
	}
	// 第一次内存,没有数据时,需要从硬盘读入数据.
	g.locker.Lock()
	g.cachedInfo = *cachedInfo
	g.locker.Unlock()
	if !reflect.DeepEqual(cachedInfo.ImportPacakgeList ,g.RequestImportList){
		g.reloadCache()
		return
	}
	kmgCache.MustMd5FileChangeCache("kmgViewResource_"+g.Name, cachedInfo.NeedCachePathList, g.reloadCache)
}

func (g *Generated) reloadCache() {
	debugBuildPath := kmgConfig.DefaultEnv().PathInProject("tmp/kmgViewResource_debug/" + g.Name)
	response := resourceBuildToDir(g.RequestImportList, debugBuildPath)
	response.NeedCachePathList = append(response.NeedCachePathList, debugBuildPath)
	kmgJson.MustWriteFileIndent(kmgConfig.DefaultEnv().PathInProject("tmp/kmgViewResource_meta/"+g.Name), response)
	g.locker.Lock()
	g.cachedInfo = response
	g.locker.Unlock()
}

/*
TODO 修复缓存无效原因1:
	* 用户修改了资源文件的配置,现在资源文件需要访问的package和上一次缓存时需要访问的package不一致.
 */