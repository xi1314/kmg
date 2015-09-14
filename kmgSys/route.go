package kmgSys

import (
	"fmt"
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgPlatform"
	"strings"
)

type RouteRule struct {
	Destination string
	Gateway     string
	Genmask     string
	Iface       string
}

// 仅支持linux ipv4
func MustGetRouteTable() []*RouteRule {
	if !kmgPlatform.IsLinux() {
		panic("[MustGetRouteTable] only support linux now")
	}
	outputS := kmgCmd.MustCombinedOutput("netstat -nr")
	lineList := strings.Split(string(outputS), "\n")
	if len(lineList) < 2 {
		panic("[getRouteTable] len(lineList)<2 " + string(outputS))
	}
	output := []*RouteRule{}
	for i, line := range lineList[2:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		part := strings.Fields(line)
		if len(part) < 4 {
			panic(fmt.Errorf("[getRouteTable] len(part)<4 lineNum:%d content:%s", i, string(outputS)))
		}
		thisRule := &RouteRule{
			Destination: part[0],
			Gateway:     part[1],
			Genmask:     part[2],
			Iface:       part[len(part)-1],
		}
		output = append(output, thisRule)
	}
	return output
}

// 仅支持linux ipv4
// 返回nil表示没找到.
func MustGetDefaultRoute() *RouteRule {
	if !kmgPlatform.IsLinux() {
		panic("[MustGetDefaultRoute] only support linux now")
	}
	for _, route := range MustGetRouteTable() {
		if route.Genmask == "0.0.0.0" {
			return route
		}
	}
	return nil
}
