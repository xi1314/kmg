package kmgSys

import (
	"github.com/bronze1man/kmgTest"
	"testing"
)

func TestMemory(t *testing.T) {
	out := `
Mem:    2099130368 1629151232  469979136   59084800   41701376  557920256
-/+ buffers/cache: 1029529600 1069600768
Swap:            0          0          0
	`
	used, total := memory(out)
	kmgTest.Equal(total, 2099130368)
	kmgTest.Equal(used, 0.49046)

	out = `
Mem:    33727324160 32708145152 1019179008          0  466333696 29958025216
-/+ buffers/cache: 2283786240 31443537920
Swap:   34347151360   47198208 34299953152
	`
	used, total = memory(out)
	kmgTest.Equal(total, 2283786240+31443537920)
	kmgTest.Equal(used, 0.067713)
}

func TestCpu(t *testing.T) {
	out := `Linux 3.13.0-55-generic (DEV-BCE) 	07/15/2015 	_x86_64_	(2 CPU)

09:16:45 AM  CPU    %usr   %nice    %sys %iowait    %irq   %soft  %steal  %guest  %gnice   %idle
09:16:45 AM  all    0.45    0.01    0.17    0.27    0.00    0.00    0.00    0.00    0.00   99.10
`
	used, count := cpu(out)
	kmgTest.Equal(used, 0.009)
	kmgTest.Equal(count, 2)
}

func TestDisk(t *testing.T) {
	out := `Filesystem     1K-blocks     Used Available Use% Mounted on
/dev/vda1       20510264  6824604  12778152  35% /
none                   4        0         4   0% /sys/fs/cgroup
udev             1014108       12   1014096   1% /dev
tmpfs             204996      436    204560   1% /run
none                5120        0      5120   0% /run/lock
none             1024964        0   1024964   0% /run/shm
none              102400        0    102400   0% /run/user
/dev/vdb1      103080224 11160216  86660796  12% /mnt
	`
	used, total := disk(out)
	kmgTest.Equal(total, 20510264)
	kmgTest.Equal(used, 0.35)
}

func TestNetwork(t *testing.T) {
	out := `eth0      Link encap:Ethernet  HWaddr 10:bf:48:4f:08:20
	  inet addr:222.197.183.79  Bcast:222.197.183.95  Mask:255.255.255.224
	  inet6 addr: fe80::12bf:48ff:fe4f:820/64 Scope:Link
	  UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1
	  RX packets:77719724 errors:0 dropped:0 overruns:0 frame:0
	  TX packets:84835236 errors:0 dropped:0 overruns:0 carrier:0
	  collisions:0 txqueuelen:1000
	  RX bytes:28383392489 (28.3 GB)  TX bytes:38417002289 (38.4 GB)
	  Interrupt:22 Memory:fba00000-fba20000
`
	rx, tx := networkRXTX(out)
	kmgTest.Equal(rx, 28383392489)
	kmgTest.Equal(tx, 38417002289)
}

func TestNetwork1(t *testing.T) {
	out := `eth0      Link encap:Ethernet  HWaddr fa:16:3e:e5:c7:08
          inet addr:192.168.0.7  Bcast:192.168.255.255  Mask:255.255.0.0
          inet6 addr: fe80::f816:3eff:fee5:c708/64 Scope:Link
          UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1
          RX packets:1380063 errors:0 dropped:0 overruns:0 frame:0
          TX packets:1188991 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:1000
          RX bytes:412780439 (412.7 MB)  TX bytes:1187905033 (1.1 GB)
`
	rx, tx := networkRXTX(out)
	kmgTest.Equal(rx, 412780439)
	kmgTest.Equal(tx, 1187905033)
}

func TestNetwork2(t *testing.T) {
	out := `eth0      Link encap:Ethernet  HWaddr bc:76:4e:1c:24:bc
          inet addr:119.9.108.209  Bcast:119.9.108.255  Mask:255.255.255.0
          inet6 addr: 2401:1800:7800:104:be76:4eff:fe1c:24bc/64 Scope:Global
          inet6 addr: fe80::be76:4eff:fe1c:24bc/64 Scope:Link
          UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1
          RX packets:2926371 errors:0 dropped:0 overruns:0 frame:0
          TX packets:3082883 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:1000
          RX bytes:2223211283 (2.2 GB)  TX bytes:2315810417 (2.3 GB)
	`
	rx, tx := networkRXTX(out)
	kmgTest.Equal(rx, 2223211283)
	kmgTest.Equal(tx, 2315810417)
}

func TestNetworkConnection(t *testing.T) {
	count := networkConnection(`
	39`)
	kmgTest.Equal(count, 39)
}

func TestIKEUserCount(t *testing.T) {
	c := ikeUserCount(`uptime: 45 minutes, since Jul 17 16:03:02 2015
worker threads: 32 total, 27 idle, working: 4/0/1/0
job queues: 0/0/0/0
jobs scheduled: 471
IKE_SAs: 79 total, 0 half-open
mallinfo: sbrk 7364608, mmap 0, used 2095152, free 5269456`)
	kmgTest.Equal(c, 79)
}
