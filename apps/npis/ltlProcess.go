package npis

import (
	"github.com/slzm40/gogate/apps/cacheq"
	"github.com/slzm40/gogate/apps/mq"
	"github.com/slzm40/gogate/controllers/elinkpsh"
	"github.com/slzm40/gogate/models/devmodels"
	"github.com/slzm40/gogate/protocol/elinkres"
	"github.com/slzm40/gomo/elink"
	"github.com/slzm40/gomo/ltl"
	"github.com/slzm40/gomo/protocol/limp"

	"github.com/astaxie/beego/logs"
)

func (this *ZbnpiApp) ProcessInSpecificCmd(srcAddr uint16, hdr *ltl.FrameHdr, cmdFormart []byte) byte {
	return 0
}
func (this *ZbnpiApp) ProcessInReadCmd(srcAddr uint16, hdr *ltl.FrameHdr, attrId []uint16) error {
	return nil
}
func (this *ZbnpiApp) ProcessInReadRspCmd(srcAddr uint16, hdr *ltl.FrameHdr, rdRspStatus []ltl.RcvReadRspStatus) error {
	itm, err := cacheq.Excute(hdr.TransSeqNum)
	if err != nil {
		return err
	}
	switch hdr.TrunkID {
	case ltl.TrunkID_GeneralBasic:
		gba := limp.BasicAttribute(rdRspStatus)
		if itm.IsLocal {
			s, ok := itm.Val.(string)
			if ok {
				devmodels.UpdateZbDeviceAndNode(s, srcAddr, 1, gba.ProductIdentifier)
				if IsNetworkSteering() {
					elinkpsh.DeviceAnnce(gba.ProductIdentifier, s, true)
				}
			}
		} else {

		}
	default:
		return nil
	}

	return nil
}
func (this *ZbnpiApp) ProcessInWriteCmd(srcAddr uint16, hdr *ltl.FrameHdr, wrwrRec []ltl.RcvWriteRec) error {
	return nil
}
func (this *ZbnpiApp) ProcessInWriteUndividedCmd(srcAddr uint16, hdr *ltl.FrameHdr, wrwrRec []ltl.RcvWriteRec) error {
	return nil
}
func (this *ZbnpiApp) ProcessInWriteRspCmd(srcAddr uint16, hdr *ltl.FrameHdr, wrStatus []ltl.WriteRspStatus) error {
	return nil
}
func (this *ZbnpiApp) ProcessInConfigReportCmd(srcAddr uint16, hdr *ltl.FrameHdr, crRec []ltl.RcvCfgReportRec) error {
	return nil
}
func (this *ZbnpiApp) ProcessInConfigReportRspCmd(srcAddr uint16, hdr *ltl.FrameHdr, crStatus []ltl.CfgReportRspStatus) error {
	return nil
}
func (this *ZbnpiApp) ProcessInReadConfigReportCmd(srcAddr uint16, hdr *ltl.FrameHdr, attrId []uint16) error {
	return nil
}
func (this *ZbnpiApp) ProcessInReadConfigReportRspCmd(srcAddr uint16, hdr *ltl.FrameHdr, rcStatus []ltl.RcvReportCfgRspStatus) error {
	return nil
}

func (this *ZbnpiApp) ProcessInReportCmd(srcAddr uint16, hdr *ltl.FrameHdr, rRec []ltl.RcvReportRec) error {
	var out []byte

	zbdnode, err := devmodels.LookupZbDeviceNodeByNN(srcAddr, hdr.NodeNo)
	if err != nil {
		return err
	}

	switch hdr.TrunkID {
	case ltl.TrunkID_MsTemperatureMeasurement, ltl.TrunkID_MsRelativeHumidity:
		out, err = limp.MsMeasureAttrReport(zbdnode.Sn, zbdnode.ProductId, int(hdr.NodeNo),
			hdr.TrunkID, rRec)
		if err != nil {
			return err
		}
	default:
		logs.Error("no fix trunkid")
		return nil
	}

	return mq.WriteCtrlData(
		elink.FormatResouce(elinkres.DevicePropertys, zbdnode.ProductId),
		elink.MethodPatch, elink.MessageTypeAnnce, out)
}

func (this *ZbnpiApp) ProcessInDefaultRsp(srcAddr uint16, hdr *ltl.FrameHdr, dfStatus *ltl.DefaultRsp) error {
	return nil
}
