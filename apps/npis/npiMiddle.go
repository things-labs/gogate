package npis

import (
	"errors"

	"github.com/thinkgos/gomo/ltl"
	"github.com/thinkgos/gomo/npi"
)

// npi与ltl的中间层
type MiddleMonitor struct {
	*npi.Monitor
	IncommingMsgPkt chan *ltl.IncomingMsgPkt
}

func NewMiddleMonitor(m *npi.Monitor) *MiddleMonitor {
	mid := &MiddleMonitor{
		IncommingMsgPkt: make(chan *ltl.IncomingMsgPkt, Incomming_msg_size_max),
		Monitor:         m,
	}

	mid.AddAsyncCbs(map[uint16]func(*npi.Npdu){
		npi.MT_AF_DATA_CONFIRM:                        Af_DataConfirm,
		npi.MT_AF_INCOMING_MSG:                        Af_IncomingMsg,
		npi.MT_ZDO_MGMT_PERMIT_JOIN_RSP:               Zdo_MgmtPermitJoinRsp,
		npi.MT_ZDO_STATE_CHANGE_IND:                   Zdo_StateChangeInd,
		npi.MT_ZDO_END_DEVICE_ANNCE_IND:               Zdo_EnddeviceAnnceInd,
		npi.MT_ZDO_LEAVE_IND:                          Zdo_LeaveInd,
		npi.MT_SYS_RESET_IND:                          Sys_ResetInd,
		npi.MT_APP_CNF_BDB_COMMISSIONING_NOTIFICATION: Appcfg_BdbCommissioningNotice,
	})

	return mid
}

func (this *MiddleMonitor) writeIncommingMsg(msg *ltl.IncomingMsgPkt) {
	select {
	case <-this.CloseChan():
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
