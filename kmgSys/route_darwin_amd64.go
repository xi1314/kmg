package kmgSys

import (
	"fmt"
	"golang.org/x/sys/unix"
	"unsafe"
	//"github.com/bronze1man/kmg/kmgDebug"
	"net"
)

type Route struct {
	Destination net.IPNet
	GateWay     net.IP
}

func GetRouteTable() (routeList []*Route, err error) {
	buf, err := unix.RouteRIB(unix.NET_RT_DUMP, 0)
	if err != nil {
		return nil, err
	}
	/*
		msgs:=parseAnyMessage(buf)
		kmgDebug.Println(msgs)
		fmt.Println(msgs)
	*/
	msgs, err := unix.ParseRoutingMessage(buf)
	if err != nil {
		return nil, err
	}
	routeList = make([]*Route, len(msgs))
	for _, msg := range msgs {
		addrList, err := unix.ParseRoutingSockaddr(msg)
		if err != nil {
			return nil, err
		}
		for _, addr := range addrList {
			fmt.Printf("%T\n", addr)
		}
		//route:=&Route{

		//}
		//r = append(r,route)
		//fmt.Printf("%T %#v %#v\n",msg,msg,msg.(*unix.RouteMessage))
	}
	return nil, nil
}

func parseAnyMessage(b []byte) (msgs []*anyMessage) {
	msgCount := 0
	for len(b) >= anyMessageLen {
		msgCount++
		any := (*anyMessage)(unsafe.Pointer(&b[0]))
		if any.Version != unix.RTM_VERSION {
			b = b[any.Msglen:]
			continue
		}
		msgs = append(msgs, any)
		b = b[any.Msglen:]
	}
	return msgs
}

type anyMessage struct {
	Msglen  uint16
	Version uint8
	Type    uint8
}

const anyMessageLen = int(unsafe.Sizeof(anyMessage{}))
