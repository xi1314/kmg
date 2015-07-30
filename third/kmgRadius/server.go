package kmgRadius

import (
	"fmt"

	. "github.com/bronze1man/kmg/kmgErr"
)

type server struct {
	mschapMap map[string]mschapStatus
	handler   Handler
}

type mschapStatus struct {
	Challenge  [16]byte
	NTResponse [24]byte
}

func (p *server) PacketHandler(request *Packet) *Packet {
	switch request.Code {
	case CodeAccessRequest:
		return p.radiusAccess(request)
	case CodeAccountingRequest:
		return p.radiusAccountingRequest(request)
	default:
		npac := request.Reply()
		LogError(fmt.Errorf("[radius.RadiusHandle] request.Code %s", request.Code.String()))
		npac.Code = CodeAccessReject
		return npac
	}
}
