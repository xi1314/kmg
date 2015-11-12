package kmgViewResource

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgCrypto"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgGoSource/kmgFormat"
	"github.com/bronze1man/kmg/kmgStrings"
	"github.com/bronze1man/kmg/third/kmgQiniu"
	"path"
	"path/filepath"
	"sort"
)

type ResourceUploadRequest struct {
	ImportPathList []string
	//ResourceList  []string //传入一堆资源目录的列表,然后传到七牛上.
	Qiniu         *kmgQiniu.Context // 如果传入qiniu相关的对象,生成代码时,会向七牛上传资源文件
	QiniuPrefix   string
	OutGoFilePath string
	Name          string //缓存和区分不同实例使用.
}

var allowResourceExt = []string{".otf", ".eot", ".svg", ".ttf", ".woff", ".woff2", ".jpg", ".jpeg", ".png", ".gif", ".ico", ".html"}

func ResourceBuild(req *ResourceUploadRequest) {
	if req.Name == "" {
		panic(`[ResourceBuild] req.Name == ""`)
	}
	tmpDirPath := kmgConfig.DefaultEnv().PathInTmp("kmgViewResource_build/" + req.Name)
	kmgFile.MustDelete(tmpDirPath)
	response := resourceBuildToDir(req.ImportPathList, tmpDirPath)
	req.Qiniu.MustUploadFromFile(tmpDirPath, req.QiniuPrefix)

	packageName := filepath.Base(filepath.Dir(req.OutGoFilePath))

	urlPrefix := req.Qiniu.GetSchemeAndDomain() + "/" + req.QiniuPrefix
	//jsUrl:=urlPrefix+"/"+response.JsFileName
	//cssUrl:=urlPrefix+"/"+response.CssFileName

	// 不可以使用 fmt.Sprintf("%#v",generated) 会导出私有变量.
	//generated:=&Generated{
	//	Name: req.Name,
	//	GeneratedJsFileUrl: jsUrl,
	//	GeneratedCssFileUrl: cssUrl,
	//	GeneratedUrlPrefix: urlPrefix,
	//	RequestImportList: req.ImportPathList,
	//}
	outGoContent := []byte(`package ` + packageName + `
import (
	"github.com/bronze1man/kmg/kmgView/kmgViewResource"
	"sync"
)
var ` + req.Name + `Once sync.Once
var ` + req.Name + `generated *kmgViewResource.Generated
func get` + req.Name + `ViewResource() *kmgViewResource.Generated{
	` + req.Name + `Once.Do(func(){
		` + req.Name + `generated = &kmgViewResource.Generated{
			Name: ` + fmt.Sprintf("%#v", req.Name) + `,
			GeneratedJsFileName: ` + fmt.Sprintf("%#v", response.JsFileName) + `,
			GeneratedCssFileName: ` + fmt.Sprintf("%#v", response.CssFileName) + `,
			GeneratedUrlPrefix: ` + fmt.Sprintf("%#v", urlPrefix) + `,
			RequestImportList: ` + fmt.Sprintf("%#v", req.ImportPathList) + `,
		}
		` + req.Name + `generated.Init()
	})
	return ` + req.Name + `generated
}
`)
	outGoContent, err := kmgFormat.Source(outGoContent)
	if err != nil {
		panic(err)
	}
	kmgFile.MustWriteFile(req.OutGoFilePath, outGoContent)
}

type resourceBuildToDirResponse struct {
	NeedCachePathList []string
	CssFileName       string
	JsFileName        string

	ImportPackageList []string //需要导入的 package 的列表, 此处对这个信息进行缓存验证
}

func resourceBuildToDir(ImportPackageList []string, tmpDirPath string) (response resourceBuildToDirResponse) {
	builder := &tBuilder{
		pkgMap: map[string]*pkg{},
	}
	for _, importPath := range ImportPackageList {
		builder.handlePkg(importPath)
	}
	for _, pkg := range builder.pkgDepOrder {
		builder.JsContent = append(builder.JsContent, pkg.JsContent...)
		builder.JsContent = append(builder.JsContent, byte('\n'))
		builder.CssContent = append(builder.CssContent, pkg.CssContent...)
		builder.CssContent = append(builder.CssContent, byte('\n'))
	}

	response.CssFileName = kmgCrypto.Md5Hex(builder.CssContent) + ".css"
	response.JsFileName = kmgCrypto.Md5Hex(builder.JsContent) + ".js"

	kmgFile.MustMkdir(tmpDirPath)
	kmgFile.MustWriteFile(filepath.Join(tmpDirPath, response.CssFileName), builder.CssContent)
	kmgFile.MustWriteFile(filepath.Join(tmpDirPath, response.JsFileName), builder.JsContent)
	for _, pkg := range builder.pkgDepOrder {
		for _, filePath := range pkg.ResourceFilePathList {
			kmgFile.MustWriteFile(filepath.Join(tmpDirPath, filepath.Base(filePath)), kmgFile.MustReadFile(filePath))
		}
	}
	for _, pkg := range builder.pkgDepOrder {
		response.NeedCachePathList = append(response.NeedCachePathList, pkg.Dirpath)
	}
	response.ImportPackageList = ImportPackageList
	return response
}

type tBuilder struct {
	pkgMap      map[string]*pkg // pkg是否访问过表.保证一个pkg在后面只出现一次
	pkgDepOrder []*pkg          // 依赖引用树,从叶子开始遍历,保证前面的不会引用后面的.
	pkgDepStack pkgStack        // 依赖循环检查堆栈,保证系统不存在依赖循环.

	JsContent         []byte
	CssContent        []byte
	CacheNeedCheckDir []string

	ResourceFileNameMap map[string]bool // 不允许资源文件的名称完全相同.
}

type pkg struct {
	PackageName    string
	Dirpath        string
	ImportPathList []string

	JsFilePathList  []string
	CssFilePathList []string

	ResourceFilePathList []string
	// 合并好的js和css的内容,
	JsContent  []byte
	CssContent []byte
}

func (b *tBuilder) handlePkg(packageName string) {
	for _, thisPkg := range b.pkgDepStack.arr {
		if thisPkg.PackageName == packageName {
			panic("[kmgViewResource] import circle " + packageName)
		}
	}
	thisPkg, ok := b.pkgMap[packageName]
	if ok {
		return
	}
	thisPkg = b.parsePkg(packageName)
	b.pkgMap[packageName] = thisPkg
	b.pkgDepStack.push(thisPkg)
	for _, importPath := range thisPkg.ImportPathList {
		b.handlePkg(importPath)
	}
	b.pkgDepStack.pop()
	b.pkgDepOrder = append(b.pkgDepOrder, thisPkg)
}

func (b *tBuilder) parsePkg(packageName string) *pkg {
	thisPkg := &pkg{
		PackageName: packageName,
	}
	thisPkg.Dirpath = path.Join(kmgConfig.DefaultEnv().GetFirstGOPATH(), "src", packageName)
	if !kmgFile.MustDirectoryExist(thisPkg.Dirpath) {
		panic("[kmgViewResource] can not found dir " + thisPkg.Dirpath)
	}
	fileList := kmgFile.MustGetAllFileOneLevel(thisPkg.Dirpath)
	for _, file := range fileList {
		ext := kmgFile.GetExt(file)
		if ext == ".js" {
			thisPkg.JsFilePathList = append(thisPkg.JsFilePathList, file)

			importPathList := parseImportPath(file, kmgFile.MustReadFile(file))
			thisPkg.ImportPathList = kmgStrings.SliceNoRepeatMerge(thisPkg.ImportPathList, importPathList)
		} else if ext == ".css" {
			thisPkg.CssFilePathList = append(thisPkg.CssFilePathList, file)

			importPathList := parseImportPath(file, kmgFile.MustReadFile(file))
			thisPkg.ImportPathList = kmgStrings.SliceNoRepeatMerge(thisPkg.ImportPathList, importPathList)
		} else if kmgStrings.IsInSlice(allowResourceExt, ext) {
			name := filepath.Base(file)
			if b.ResourceFileNameMap[name] {
				panic("[kmgViewResource] resource file name " + name + " repeat path " + file)
			}
			thisPkg.ResourceFilePathList = append(thisPkg.ResourceFilePathList, file)
		}
	}
	sort.Strings(thisPkg.JsFilePathList)
	sort.Strings(thisPkg.CssFilePathList)

	for _, file := range thisPkg.JsFilePathList {
		// 这个泄漏信息比较严重.暂时关掉吧.
		//thisPkg.JsContent = append(thisPkg.JsContent, []byte("\n/* "+file+" */\n\n")...)
		thisPkg.JsContent = append(thisPkg.JsContent, kmgFile.MustReadFile(file)...)
		thisPkg.JsContent = append(thisPkg.JsContent, byte('\n'))
	}

	for _, file := range thisPkg.CssFilePathList {
		//thisPkg.CssContent = append(thisPkg.CssContent, []byte("\n/* "+file+"*/\n\n")...)
		thisPkg.CssContent = append(thisPkg.CssContent, kmgFile.MustReadFile(file)...)
		thisPkg.CssContent = append(thisPkg.CssContent, byte('\n'))
	}

	return thisPkg
}

type pkgStack struct {
	arr []*pkg
	pos int
}

func (stack *pkgStack) push(p *pkg) {
	stack.arr = append(stack.arr, p)
	stack.pos++
}

func (stack *pkgStack) pop() *pkg {
	if stack.pos == 0 {
		panic("[pkgStack.pop] stack.pos==0")
	}
	stack.pos--
	p := stack.arr[stack.pos]
	stack.arr = stack.arr[:stack.pos]
	return p
}
