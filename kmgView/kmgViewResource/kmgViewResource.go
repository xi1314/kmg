package kmgViewResource
import (
	"github.com/bronze1man/kmg/third/kmgQiniu"
	"path"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgFile"
	"path/filepath"
	"sort"
	"github.com/bronze1man/kmg/kmgStrings"
	"github.com/bronze1man/kmg/kmgCrypto"
	"fmt"
)

type ResourceUploadRequest struct {
	EnterPointPackageName string
	//ResourceList  []string //传入一堆资源目录的列表,然后传到七牛上.
	Qiniu         *kmgQiniu.Context
	QiniuPrefix   string
	OutGoFilePath string
	FuncPrefix    string
}

func ResourceBuild(req *ResourceUploadRequest){
	builder:=&tBuilder{
		pkgMap: map[string]*pkg{},
	}
	builder.handlePkg(req.EnterPointPackageName)
	for _,pkg:=range builder.pkgDepOrder{
		builder.JsContent = append(builder.JsContent,pkg.JsContent...)
		builder.JsContent = append(builder.JsContent,byte('\n'))
		builder.CssContent = append(builder.CssContent,pkg.CssContent...)
		builder.CssContent = append(builder.CssContent,byte('\n'))
	}

	cssFileName := kmgCrypto.Md5Hex(builder.CssContent) + ".css"
	jsFileName := kmgCrypto.Md5Hex(builder.JsContent) + ".js"
	req.Qiniu.MustUploadFromBytes(req.QiniuPrefix+"/"+cssFileName, builder.CssContent)
	req.Qiniu.MustUploadFromBytes(req.QiniuPrefix+"/"+jsFileName, builder.JsContent)

	packageName := filepath.Base(filepath.Dir(req.OutGoFilePath))
	if req.FuncPrefix == "" {
		req.FuncPrefix = "getManaged"
	}
	outGoContent := []byte(`package ` + packageName + `
func ` + req.FuncPrefix + `JsUrl()string{
	return ` + fmt.Sprintf("%#v", req.Qiniu.GetSchemeAndDomain()+"/"+req.QiniuPrefix+"/"+jsFileName) + `
}
func ` + req.FuncPrefix + `CssUrl()string{
	return ` + fmt.Sprintf("%#v", req.Qiniu.GetSchemeAndDomain()+"/"+req.QiniuPrefix+"/"+cssFileName) + `
}
`)
	kmgFile.MustWriteFile(req.OutGoFilePath, outGoContent)
}

type tBuilder struct{
	pkgMap map[string]*pkg // pkg是否访问过表.保证一个pkg在后面只出现一次
	pkgDepOrder []*pkg // 依赖引用树,从叶子开始遍历,保证前面的不会引用后面的.
	pkgDepStack pkgStack // 依赖循环检查堆栈,保证系统不存在依赖循环.

	JsContent []byte
	CssContent []byte
	CacheNeedCheckDir []string
}

type pkg struct{
	PackageName string
	ImportPathList []string

	JsFilePathList []string
	CssFilePathList []string

	// 合并好的js和css的内容,
	JsContent []byte
	CssContent []byte
}

func (b *tBuilder)handlePkg(packageName string){
	for _,thisPkg:=range b.pkgDepStack.arr{
		if thisPkg.PackageName==packageName{
			panic("[kmgViewResource] import circle "+packageName)
		}
	}
	thisPkg,ok:=b.pkgMap[packageName]
	if ok{
		return thisPkg
	}
	thisPkg = b.parsePkg(packageName)
	b.pkgMap[packageName] = thisPkg
	b.pkgDepStack.push(thisPkg)
	for _,importPath:=range thisPkg{
		b.handlePkg(importPath)
	}
	b.pkgDepStack.pop()
	b.pkgDepOrder = append(b.pkgDepOrder,thisPkg)
}

func (b *tBuilder) parsePkg(packageName string)*pkg{
	thisPkg:=&pkg{
		packageName: packageName,
	}
	dirPath:=path.Join(kmgConfig.DefaultEnv().GetGOROOT(),"src",packageName)
	if !kmgFile.MustDirectoryExist(dirPath){
		panic("[kmgViewResource] can not found dir "+dirPath)
	}
	fileList:=kmgFile.MustGetAllFiles(dirPath)
	for _,file:=range fileList{
		ext:= kmgFile.GetExt(file)
		if ext==".js" {
			thisPkg.JsFilePathList = append(thisPkg.JsFilePathList,file)

			importPathList:=parseImportPath(file,kmgFile.MustReadFile(file))
			thisPkg.ImportPathList = kmgStrings.SliceNoRepeatMerge(thisPkg.ImportPathList,importPathList)
		}else if ext==".css"{
			thisPkg.CssFilePathList = append(thisPkg.CssFilePathList,file)

			importPathList:=parseImportPath(file,kmgFile.MustReadFile(file))
			thisPkg.ImportPathList = kmgStrings.SliceNoRepeatMerge(thisPkg.ImportPathList,importPathList)
		}
	}
	sort.Strings(thisPkg.JsFilePathList)
	sort.Strings(thisPkg.CssFilePathList)

	for _,file:=range thisPkg.JsFilePathList{
		thisPkg.JsContent = append(thisPkg.JsContent,kmgFile.MustReadFile(file)...)
		thisPkg.JsContent = append(thisPkg.JsContent,byte('\n'))
	}

	for _,file:=range thisPkg.CssFilePathList{
		thisPkg.CssContent = append(thisPkg.CssContent,kmgFile.MustReadFile(file)...)
		thisPkg.CssContent = append(thisPkg.CssContent,byte('\n'))
	}

	return thisPkg
}

type pkgStack struct{
	arr []*pkg
	pos int
}

func (stack *pkgStack) push(p *pkg){
	stack.arr = append(stack.arr,p)
	stack.pos++
}

func (stack *pkgStack) pop() *pkg{
	if stack.pos==0{
		panic("[pkgStack.pop] stack.pos==0")
	}
	stack.pos--
	p:=stack.arr[stack.pos]
	stack.arr = stack.arr[:stack.pos]
	return p
}
