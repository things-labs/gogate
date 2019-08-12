package npis

import (
	"github.com/thinkgos/gogate/controllers/elinkpsh"
	"github.com/thinkgos/gogate/models"
	"github.com/thinkgos/gomo/ltl"
	"github.com/thinkgos/gomo/npi"
	"github.com/thinkgos/memlog"
)

/** AF **/
func Af_DataConfirm(pdu *npi.Npdu) {
	_, err := npi.Af_DataConfirmParse(pdu)
	if err != nil {
		memlog.Error("af data confirm: ", err)
		return
	}

	memlog.Debug("af data confirm success")
}

func Af_IncomingMsg(pdu *npi.Npdu) {
	o, err := npi.Af_IncomingMsgParse(pdu)
	if err != nil {
		memlog.Error(err)
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
		memlog.Error("zdo permit join rsp: %s", err)
		return
	}

	memlog.Debug("zdo permit join rsp: %#v", o)
}

func Zdo_StateChangeInd(pdu *npi.Npdu) {
	o, err := npi.Zdo_StateChangeIndParse(pdu)
	if err != nil {
		memlog.Error("zdo state indicate: %s", err)
		return
	}
	memlog.Debug("zdo state indicate: %#v", o)
	if o == 0x09 {
		ZbApps.isNetworkFormation = true
	} else {
		ZbApps.isNetworkFormation = false
	}
}

func Zdo_EnddeviceAnnceInd(pdu *npi.Npdu) {
	o, err := npi.Zdo_EnddeviceAnnceIndParse(pdu)
	if err != nil {
		memlog.Error("enddevice annce: %s", err)
		return
	}
	memlog.Debug("enddevice annce: [%s - %d]", models.ToHexString(o.IeeeAddr), o.NwkAddr)

	ZbApps.SendReadReqBasic(o.NwkAddr, deviceAnnce{models.ToHexString(o.IeeeAddr), o.NwkAddr})
}

func Zdo_LeaveInd(pdu *npi.Npdu) {
	o, err := npi.Zdo_LeaveIndParse(pdu)
	if err != nil {
		memlog.Error("leave indicate: %s", err)
		return
	}
	sn := models.ToHexString(o.ExtAddr)
	dev, err := models.LookupZbDeviceByIeeeAddr(sn)
	if err != nil {
		return
	}

	err = dev.DeleteZbDeveiceAndNode()
	if err != nil {
		memlog.Error("leave indicate: %s", err)
		return
	}

	elinkpsh.DeviceAnnce(dev.GetProductID(), dev.GetSn(), false)
	memlog.Debug("levae indicate: [%s - %d]", dev.GetSn(), dev.GetNwkAddr())
}

/** SYS **/
func Sys_ResetInd(pdu *npi.Npdu) {
	o, err := npi.Sys_ResetIndParse(pdu)
	if err != nil {
		memlog.Error("reset indicate: %s", err)
		return
	}

	memlog.Debug("reset indicate: %#v", o)
}

/** APPCFG **/
func Appcfg_BdbCommissioningNotice(pdu *npi.Npdu) {
	o, err := npi.Appcfg_BdbCommissioningNoticeParse(pdu)
	if err != nil {
		memlog.Error(err)
		return
	}

	memlog.Debug("bdb notice: %#v", o)
}
