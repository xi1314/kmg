package kmgSys

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgMath"
	"github.com/bronze1man/kmg/kmgStrconv"
	"strings"
)

func UsagePreInstall() {
	if kmgCmd.Exist("mpstat") && kmgCmd.Exist("netstat") {
		return
	}
	kmgCmd.MustRun("sudo apt-get update")
	kmgCmd.MustRun("sudo apt-get install -y sysstat")
}

//byte
func Memory() (used float64, total int) {
	return memory(string(kmgCmd.MustRunAndReturnOutput("free -b")))
}

func memory(output string) (used float64, total int) {
	o := strings.Split(output, "\n")
	o0 := strings.TrimPrefix(o[2], "-/+ buffers/cache:")
	o1 := strings.Split(o0, " ")
	usedByte := 0
	free := 0
	for _, s := range o1 {
		if s == "" {
			continue
		}
		if usedByte == 0 {
			usedByte = kmgStrconv.AtoIDefault0(s)
			continue
		}
		free = kmgStrconv.AtoIDefault0(s)
		break
	}
	total = usedByte + free
	used = float64(usedByte) / float64(total)
	used = kmgMath.Float64RoundToRelativePrec(used, 4)
	return
}

func Cpu() (used float64, numOfCore int) {
	return cpu(string(kmgCmd.MustRunAndReturnOutput("mpstat")))
}

func cpu(output string) (used float64, numOfCore int) {
	o := strings.Split(output, "\n")
	o0 := strings.TrimSuffix(o[0], " CPU)")
	ol0 := strings.Split(o0, "(")
	numOfCore = kmgStrconv.AtoIDefault0(ol0[len(ol0)-1])
	o1 := strings.Split(o[3], " ")
	_used := kmgStrconv.ParseFloat64Default0(o1[len(o1)-1])
	used = (float64(100) - _used) / float64(100)
	used = kmgMath.Float64RoundToRelativePrec(used, 4)
	return
}

//只返回 / 挂载的磁盘空间
//total 1024byte 的默认单位，没有选项改成byte
func Disk() (used float64, total int) {
	return disk(string(kmgCmd.MustRunAndReturnOutput("df")))
}

func disk(output string) (used float64, total int) {
	o := strings.Split(output, "\n")
	o0 := strings.Split(o[1], " ")
	after := 0
	for i, s := range o0 {
		if i == 0 {
			continue
		}
		if s == "" {
			continue
		}
		if total == 0 {
			total = kmgStrconv.AtoIDefault0(s)
			continue
		}
		after++
		if after == 3 {
			p := kmgStrconv.ParseFloat64Default0(strings.TrimSuffix(s, "%"))
			used = kmgMath.Float64RoundToRelativePrec(p/float64(100), 4)
			break
		}
	}
	return
}

//byte
func NetworkRXTX(deviceName string) (rx int, tx int) {
	return networkRXTX(string(kmgCmd.MustRunAndReturnOutput("ifconfig " + deviceName)))
}

func networkRXTX(output string) (rx int, tx int) {
	al := strings.Split(output, "\n")
	sl := ""
	for _, v := range al {
		v = strings.TrimSpace(v)
		v = strings.TrimPrefix(v, "\t")
		if strings.HasPrefix(v, "RX bytes") {
			sl = v
			break
		}
	}
	l := strings.Split(sl, "  ")
	_l := []string{}
	for _, s := range l {
		if s == "" || s == "\t" {
			continue
		}
		_l = append(_l, s)
	}
	rx = kmgStrconv.AtoIDefault0(strings.TrimPrefix(strings.TrimSuffix(strings.Split(_l[0], "(")[0], " "), "RX bytes:"))
	tx = kmgStrconv.AtoIDefault0(strings.TrimPrefix(strings.TrimSuffix(strings.Split(_l[1], "(")[0], " "), "TX bytes:"))
	return
}

func NetworkConnection() (connectionCount int) {
	return networkConnection(string(kmgCmd.CmdSlice([]string{"bash", "-c", "netstat -na | grep ESTABLISHED | wc -l"}).MustRunAndReturnOutput()))
}

func networkConnection(output string) (connectionCount int) {
	output = strings.TrimSpace(output)
	return kmgStrconv.AtoIDefault0(output)
}

func IKEUserCount() int {
	if !kmgCmd.Exist("swanctl") {
		return 0
	}
	return ikeUserCount(string(kmgCmd.MustRunAndReturnOutput("swanctl -S")))
}

func ikeUserCount(output string) int {
	lines := strings.Split(output, "\n")
	c := ""
	for _, line := range lines {
		if strings.HasPrefix(line, "IKE_SAs") {
			_line := strings.Split(line, "total")
			c = _line[0]
			c = strings.Trim(c, "IKE_SAs:")
			c = strings.TrimSpace(c)
			break
		}
	}
	return kmgStrconv.AtoIDefault0(c)
}
