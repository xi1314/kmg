package kmgThirdCloud

import "github.com/bronze1man/kmg/kmgTest"

func ExampleAliyunSDK() {
	sdk := NewAliyunSDK("AccessKeyId", "AccessKeySecret", "SecurityGroupId")
	sdk.FlavorName = "ecs.t1.small"
	sdk.InstanceName = "Hello-Aliyun"
	sdk.InstancePaidType = AliyunPaidTypePost
	sdk.InstancePassword = "" //Enter your instance password
	ip := sdk.CreateInstance()
	instance, exist := sdk.ListAllRunningInstance()[ip]
	if !exist {
		panic("Rackspace CreateInstance Failed")
	}
	kmgTest.Equal(instance.Name, "Hello-Aliyun")
	sdk.RenameInstanceByIp("Bye-Aliyun", ip)
	instance, exist = sdk.ListAllRunningInstance()[ip]
	if !exist {
		panic("Aliyun Rename Failed")
	}
	kmgTest.Equal(instance.Name, "Bye-Aliyun")
	sdk.DeleteInstance(ip)
}

func ExampleRackspace() {
	sdk := NewRackspaceSDK("UserName", "ApiKey", "SSHKeyName")
	ip := sdk.CreateInstance()
	instance, exist := sdk.ListAllRunningInstance()[ip]
	if !exist {
		panic("Rackspace CreateInstance Failed")
	}
	kmgTest.Equal(instance.Name, "Hello-Rackspace")
	sdk.RenameInstanceByIp("Bye-Rackspace", ip)
	instance, exist = sdk.ListAllRunningInstance()[ip]
	if !exist {
		panic("Rackspace Rename Failed")
	}
	kmgTest.Equal(instance.Name, "Bye-Rackspace")
	sdk.DeleteInstance(ip)
}

func ExampleRackspaceForTest() {
	sdk := NewRackspaceSDK("UserName", "ApiKey", "SSHKeyName")
	sdk.FlavorName = "512MB Standard Instance"
	sdk.InstanceName = "Hello-Rackspace"
	ip := sdk.CreateInstance()
	instance, exist := sdk.ListAllRunningInstance()[ip]
	if !exist {
		panic("Rackspace CreateInstance Failed")
	}
	kmgTest.Equal(instance.Name, "Hello-Rackspace")
	sdk.RenameInstanceByIp("Bye-Rackspace", ip)
	instance, exist = sdk.ListAllRunningInstance()[ip]
	if !exist {
		panic("Rackspace Rename Failed")
	}
	kmgTest.Equal(instance.Name, "Bye-Rackspace")
	sdk.DeleteInstance(ip)
}
