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

	logs.Debug("%#v", o)

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
		logs.Error(err)
		return
	}

	logs.Debug("%#v", o)
}

func Zdo_StateChangeInd(pdu *npi.Npdu) {
	o, err := npi.Zdo_StateChangeInd(pdu)
	if err != nil {
		logs.Error(err)
		return
	}

	logs.Debug("%#v", o)
}

func Zdo_enddeviceAnnceInd(pdu *npi.Npdu) {
	o, err := npi.Zdo_EnddeviceAnnceIndParse(pdu)
	if err != nil {
		logs.Error(err)
		return
	}

	logs.Debug("%#v", o)
}

func Zdo_LeaveInd(pdu *npi.Npdu) {
	o, err := npi.Zdo_LeaveIndParse(pdu)
	if err != nil {
		logs.Error(err)
		return
	}

	logs.Debug("%#v", o)
}

/** SYS **/
func Sys_ResetInd(pdu *npi.Npdu) {
	o, err := npi.Sys_ResetIndParse(pdu)
	if err != nil {
		logs.Error(err)
		return
	}

	logs.Debug("%#v", o)
}

/** APPCFG **/
func Appcfg_BdbCommissioningNotice(pdu *npi.Npdu) {
	o, err := npi.Appcfg_BdbCommissioningNoticeParse(pdu)
	if err != nil {
		logs.Error(err)
		return
	}

	logs.Debug("%#v", o)
}
