package npis

import (
	"github.com/slzm40/gogate/apps/mq"
	"github.com/slzm40/gomo/elink"
	"github.com/slzm40/gomo/ltl"
	"github.com/slzm40/gomo/protocol/elinkch/ctrl"
	"github.com/slzm40/gomo/protocol/elinkres"
	"github.com/slzm40/gomo/protocol/limp"

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

type DevAttr struct {
	ctrl.BasePushData
	Payload struct {
		ctrl.BaseNodePayload
		Data interface{}
	} `json:"payload"`
}

func (this *ZbnpiApp) ProcessInReportCmd(srcAddr uint16, hdr *ltl.FrameHdr, rRec []ltl.RcvReportRec) error {
	//var err error
	var out []byte

	switch hdr.TrunkID {
	case ltl.TrunkID_MsTemperatureMeasurement:
		mstemp, err := limp.MsMeasureAttribute(ltl.TrunkID_MsTemperatureMeasurement, rRec)
		if err != nil {
			return err
		}
		in := DevAttr{
			Payload: struct {
				ctrl.BaseNodePayload
				Data interface{}
			}{
				BaseNodePayload: ctrl.BaseNodePayload{
					ProductID: 20000,
					Sn:        "建一个模拟的",
					NodeNo:    1,
				},
				Data: mstemp,
			},
		}
		out, err = jsoniter.Marshal(in)
		if err != nil {
			return err
		}
	case ltl.TrunkID_MsRelativeHumidity:
		mstemp, err := limp.MsMeasureAttribute(ltl.TrunkID_MsRelativeHumidity, rRec)
		if err != nil {
			return err
		}
		in := DevAttr{
			Payload: struct {
				ctrl.BaseNodePayload
				Data interface{}
			}{
				BaseNodePayload: ctrl.BaseNodePayload{
					ProductID: 20000,
					Sn:        "建一个模拟的",
					NodeNo:    2,
				},
				Data: mstemp,
			},
		}
		out, err = jsoniter.Marshal(in)
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
