使用说明:
1.写demo，eg:文件位置：src/INVE/Test/demo.go(这个名字无所谓，随便取)
type Project struct {
}
func (p *Project) GetProject()*INVE.Project{
	return INVE.GetProjectById(3)
}

2.生成代码：eg:
	kmgRpcSwift.MustGenerateCodeWithCache(&kmgRpcSwift.GenerateRequest{
		ObjectPkgPath:   "INVE/Test",
		ObjectName:      "Project",       //需要rpc通信的类名，可以调用该类的方法，变量。。。
		ObjectIsPointer: true,
		OutFilePath:     "src/INVE/Test/generated.swift", //生成位置，以及文件名
		OutClassName:    "Demo",
		OutProjectName:  "TestRpc",
		NeedSource:      true,								//第一次使用时加上这一句，会生成几个用于配置项目的文件
	})
	
3.拷贝文件
如果是第一次使用，则会生成xxx.swift,yyy-Bridging-Header.h,NSDataCompression.h,NSDataCompression.m,info.plist以及podfile文件，
将最后一个丢到ios项目根目录，其他扔到根目录下的项目名文件夹内，提示是否覆盖覆盖就行了。

4.在根目录执行 pod install ，完成后记得通过 .xworkspace 来打开。

5.打开项目后将刚才的xxx.swift,yyy-Bridging-Header.h,NSDataCompression.h,NSDataCompression.m加入项目（xcode蛋疼的文件管理），然
后在 Linked Frameworks and Libraries 中添加libz.1.2.5.tdb依赖。

完成配置，可以使用了。下面是介绍使用

eg:在ViewController的viewDidiLoad中加入下面的代码：
let client = Demo.ConfigDefaultClient("http://localhost:34895", pskStr: "abc")
print(client.GetProject())

