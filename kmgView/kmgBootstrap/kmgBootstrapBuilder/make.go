package kmgBootstrapBuilder
import "github.com/bronze1man/kmg/kmgView/kmgViewResource"

func Build(){
	kmgViewResource.ResourceBuild(&kmgViewResource.ResourceUploadRequest{
		ImportPathList: []string{
			"src/github.com/bronze1man/kmg/kmgView/kmgBootstrap/Resource/Bootstrap",
		},
		Qiniu: getKmgToolQiniu(),
		QiniuPrefix: "kmgBootstrap",
		OutGoFilePath: "src/github.com/bronze1man/kmg/kmgView/kmgBootstrap/generated_BootstrapResource.go",
		FuncPrefix: "getBootstrap",
	})
}