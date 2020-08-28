package npis

import (
	"errors"

	"github.com/thinkgos/gogate/pkg/ltl"
	"github.com/thinkgos/gogate/pkg/npi"
)

// npi与ltl的中间层
type MiddleMonitor struct {
	*npi.Monitor
	IncommingMsgPkt chan *ltl.IncomingMsgPkt
}

func NewMiddleMonitor(m *npi.Monitor) *MiddleMonitor {
	return &MiddleMonitor{
		IncommingMsgPkt: make(chan *ltl.IncomingMsgPkt, Incomming_msg_size_max),
		Monitor:         m,
	}
}

func (this *MiddleMonitor) writeIncommingMsg(msg *ltl.IncomingMsgPkt) {
	select {
	case <-this.Context().Done():
	case this.IncommingMsgPkt <- msg:
	}
}

func (this *MiddleMonitor) WriteMsg(DstAddr uint16, Data []byte) error {
	ok, err := this.Af_DataReq(DstAddr, 0xabcd, 0xaa, 0xaa, 0, 0, 0x1e, Data)
	if err != nil || !ok {
		return errors.New("response faield")
	}

	return nil
}

func (this *MiddleMonitor) IncommingMsg() <-chan *ltl.IncomingMsgPkt {
	return this.IncommingMsgPkt
}
