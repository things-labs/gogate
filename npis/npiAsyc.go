package npis

import (
	"github.com/astaxie/beego/logs"
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

func Af_IncomingMsgParse(pdu *npi.Npdu) {
	o, err := npi.Af_IncomingMsgParse(pdu)
	if err != nil {
		logs.Error(err)
		return
	}

	logs.Debug("%#v", o)
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
}

func Zdo_LeaveInd(pdu *npi.Npdu) {
	o, err := npi.Zdo_LeaveIndParse(pdu)
	if err != nil {
		logs.Error("levae indicate: %s", err)
		return
	}

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
