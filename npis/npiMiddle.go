package npis

import (
	"errors"

	"github.com/slzm40/gomo/ltl"
	"github.com/slzm40/gomo/npi"
)

// npi与ltl的中间层
type MiddleMonitor struct {
	*npi.Monitor
	IncommingMsgPkt chan *ltl.MoIncomingMsgPkt
}

func (this *MiddleMonitor) WriteMsg(DstAddr uint16, Data []byte) error {
	ok, err := this.Af_DataReq(DstAddr, 0xabcd, 0xaa, 0xaa, 0, 0, 0x1e, Data)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("response faield")
	}

	return nil
}

func (this *MiddleMonitor) IncommingMsg() <-chan *ltl.MoIncomingMsgPkt {
	return this.IncommingMsgPkt
}
