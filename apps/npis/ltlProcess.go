package npis

import (
	"github.com/slzm40/gogate/apps/lint"
	"github.com/slzm40/gogate/apps/mq"
	"github.com/slzm40/gogate/controllers/elinkres"
	"github.com/slzm40/gomo/elink"
	"github.com/slzm40/gomo/ltl"

	"github.com/json-iterator/go"
)

func (this *ZbnpiApp) ProcessInSpecificCmd(srcAddr uint16, hdr *ltl.FrameHdr, cmdFormart []byte) byte {
	return 0
}
func (this *ZbnpiApp) ProcessInReadCmd(srcAddr uint16, hdr *ltl.FrameHdr, attrId []uint16) error {
	return nil
}
func (this *ZbnpiApp) ProcessInReadRspCmd(srcAddr uint16, hdr *ltl.FrameHdr, rdRspStatus []ltl.RcvReadRspStatus) error {
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
	//var err error
	var out []byte

	switch hdr.TrunkID {
	case ltl.TrunkID_MsTemperatureMeasurement, ltl.TrunkID_MsRelativeHumidity:
		mstemp, err := lint.MsMeasureAttribute(rRec)
		if err != nil {
			return err
		}
		out, err = jsoniter.Marshal(mstemp)
		if err != nil {
			return err
		}
	}

	res := elink.FormatResouce(elinkres.DevicePropertys, 20000)
	return mq.WritePublishChData(res, elink.MethodPatch, elink.MessageTypeAnnce, out)
}
func (this *ZbnpiApp) ProcessInDefaultRsp(srcAddr uint16, hdr *ltl.FrameHdr, dfStatus *ltl.DefaultRsp) error {
	return nil
}
