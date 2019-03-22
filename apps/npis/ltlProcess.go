package npis

import (
	"github.com/astaxie/beego/logs"
	"github.com/slzm40/gomo/ltl"
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
	logs.Debug("srcaddress: 0x%04x, receive message: %#v", srcAddr, rRec)
	return nil
}
func (this *ZbnpiApp) ProcessInDefaultRsp(srcAddr uint16, hdr *ltl.FrameHdr, dfStatus *ltl.DefaultRsp) error {
	return nil
}
