package kmgSys

import (
	"github.com/bronze1man/kmg/kmgDebug"
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestExtractProcessList(t *testing.T) {
	pList := ExtractProcessListFromString(` 5393 /usr/bin/InlProxy vpnAccount -redisAddr=127.0.0.1:30001
 5424 /usr/bin/InlProxy AllInOneServerAV2 -ServerBAddr=119.9.108.86:80 -block
 5455 /usr/bin/InlProxy Tunnel.AesCtrEnd -l :30001 -n 127.0.0.1:6379 -k gD7wcBEOsH728sKkPx4fcjDiJnV
 5590 /usr/bin/InlProxy vpninfoserver -l :20004 -cmsServerAddr=127.0.0.1:20008 -vpnServerAddrList=120.24.93.239,120.24.93.239,120.24.93.239,120.24.93.239,120.24.95.104,120.24.95.104,120.24.95.104
 6625 InlProxy vpninfoserver -l :20006 -cmsServerAddr=127.0.0.1:20008 -ProjectName=snail -JpushProd
 7254 InlProxy fcFrontSnail -http=:20008
27939 /usr/bin/InlProxy DistributedLog.Explorer -MasterAddress=http://120.25.229.214:2756 -PingAddress=119.9.108.86 -WhoAmIUrl=http://120.25.229.214:29434/?n=fastCat.distributedLog.httpApi.WhoAmI`)
	kmgDebug.Println(pList)
	kmgTest.Equal(len(pList), 7)
	kmgTest.Equal(pList[1].Id, 5424)
	kmgTest.Equal(pList[5].Command, "InlProxy fcFrontSnail -http=:20008")
}
