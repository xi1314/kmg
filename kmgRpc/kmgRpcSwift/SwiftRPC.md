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
		ObjectName:      "Project",       					//需要rpc通信的类名，可以调用该类的方法，变量。。。
		ObjectIsPointer: true,
		OutFilePath:     "client/INVE/INVE/demo.swift", 	    //生成位置，以及文件名,请填到项目路径下的项目文件夹中 比如:client/INVE/INVE/demo.swift
		OutClassName:    "Demo",
		OutProjectName:  "INVE",							// 项目名
		NeedSource:      true,								//第一次使用时加上这一句，会生成几个用于配置项目的文件,之后可以去掉可以不去掉，去掉感觉会快些，而且不会在生成那些需要的依赖
	})
	

3.在根目录执行 pod install ，完成后记得通过 .xworkspace 来打开。

4.在Linked Frameworks and Libraries 中添加libz.1.2.5.tdb依赖。

完成配置，可以使用了。下面是使用介绍

eg:在ViewController的viewDidiLoad中加入下面的代码：
let client = Demo.ConfigDefaultClient("http://localhost:34895", pskStr: "abc")
print(client.GetProject())

