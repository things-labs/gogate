package npis

import (
	"github.com/astaxie/beego/logs"
	"github.com/slzm40/gogate/apps/cacheq"
	"github.com/slzm40/gogate/models/devmodels"
	"github.com/slzm40/gomo/ltl"
	"github.com/slzm40/gomo/npi"
)

/** AF **/
func Af_DataConfirm(pdu *npi.Npdu) {
	o, err := npi.Af_DataConfirmParse(pdu)
	if err != nil {
		logs.Error(err)
		return
	}

	logs.Debug("af data confirm: %#v", o)
}

func Af_IncomingMsg(pdu *npi.Npdu) {
	o, err := npi.Af_IncomingMsgParse(pdu)
	if err != nil {
		logs.Error(err)
		return
	}

	ZbApps.writeIncommingMsg(&ltl.IncomingMsgPkt{
		IsBroadCast: o.WasBroadcast > 0,
		SrcAddr:     o.SrcAddr,
		ApduData:    o.Data,
	})
}

/** ZDO **/
func Zdo_MgmtPermitJoinRsp(pdu *npi.Npdu) {
	o, err := npi.Zdo_MgmtPermitJoinRspParse(pdu)
	if err != nil {
		logs.Error("zdo permit join rsp: %s", err)
		return
	}

	logs.Debug("zdo permit join rsp: %#v", o)
}

func Zdo_StateChangeInd(pdu *npi.Npdu) {
	o, err := npi.Zdo_StateChangeIndParse(pdu)
	if err != nil {
		logs.Error("zdo state indicate: %s", err)
		return
	}
	if o == 0x09 {
		ZbApps.isNetworkFormation = true
	} else {
		ZbApps.isNetworkFormation = false
	}
	logs.Debug("zdo state indicate: %#v", o)
}

func Zdo_EnddeviceAnnceInd(pdu *npi.Npdu) {
	o, err := npi.Zdo_EnddeviceAnnceIndParse(pdu)
	if err != nil {
		logs.Error("enddevice annce: %s", err)
		return
	}

	logs.Debug("enddevice annce: %#v", o)
	id, err := cacheq.AllocID()
	if err != nil {
		return
	}

	if ZbApps.SendReadReq(o.NwkAddr, ltl.TrunkID_GeneralBasic, ltl.NodeNumRetained,
		id, []uint16{0, 1, 2, 3, 4, 5, 6, 7}) != nil {
		cacheq.FreeID(id)
		return
	}
	cacheq.Hang(id, &cacheq.CacheqItem{
		IsLocal: true,
		Cb:      cb,
		Val:     o.IeeeAddr})

}

func cb(ci *cacheq.CacheqItem) error {
	return nil
}

func Zdo_LeaveInd(pdu *npi.Npdu) {
	o, err := npi.Zdo_LeaveIndParse(pdu)
	if err != nil {
		logs.Error("leave indicate: %s", err)
		return
	}
	devmodels.DeleteZbDeveiceAndNode(devmodels.ToHexString(o.ExtAddr))
	logs.Debug("levae indicate: %#v", o)
}

/** SYS **/
func Sys_ResetInd(pdu *npi.Npdu) {
	o, err := npi.Sys_ResetIndParse(pdu)
	if err != nil {
		logs.Error("reset indicate: %s", err)
		return
	}

	logs.Debug("reset indicate: %#v", o)
}

/** APPCFG **/
func Appcfg_BdbCommissioningNotice(pdu *npi.Npdu) {
	o, err := npi.Appcfg_BdbCommissioningNoticeParse(pdu)
	if err != nil {
		logs.Error(err)
		return
	}

	logs.Debug("bdb notice: %#v", o)
}
