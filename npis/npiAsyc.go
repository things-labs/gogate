package npis

import (
	"github.com/astaxie/beego/logs"
	"github.com/slzm40/gomo/npi"
)

func Zdo_enddeviceAnnceInd(pdu *npi.Npi_pdu) {
	if annce, err := pdu.Zdo_EnddeviceAnnceIndParse(); err != nil {
		logs.Error(err)
	} else {
		logs.Debug("srcAddr: 0x%04x,nwkAddr: 0x%04x,ieeeAddr: 0x%016x,capabilities:0x%02x", annce.SrcAddr, annce.NwkAddr, annce.IeeeAddr, annce.Capabilities)
	}
}

func Sys_ResetInd(pdu *npi.Npi_pdu) {
	if ind, err := pdu.Sys_ResetIndParse(); err != nil {
		logs.Error(err)
	} else {
		logs.Debug(ind.Reason, ind.TransportRev, ind.ProductId, ind.MajorRel, ind.MinorRel, ind.HwRev)
	}
}

func init() {
	npi.RegisterAsyncCb(npi.MT_ZDO_END_DEVICE_ANNCE_IND, Zdo_enddeviceAnnceInd)
	npi.RegisterAsyncCb(npi.MT_SYS_RESET_IND, Sys_ResetInd)
}
