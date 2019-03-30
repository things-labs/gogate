package npis

import (
	"errors"

	jsoniter "github.com/json-iterator/go"
	"github.com/slzm40/gogate/apps/mq"
	"github.com/slzm40/gogate/controllers/elinkpsh"
	"github.com/slzm40/gogate/models/devmodels"
	"github.com/slzm40/gogate/protocol/elinkres"
	"github.com/slzm40/gogate/protocol/elmodels"
	"github.com/slzm40/gomo/elink"
	"github.com/slzm40/gomo/ltl"
	"github.com/slzm40/gomo/protocol/elinkch/ctrl"
	"github.com/slzm40/gomo/protocol/limp"

	"github.com/astaxie/beego/logs"
)

func (this *ZbnpiApp) ProInSpecificCmd(srcAddr uint16, hdr *ltl.FrameHdr, cmdFormart []byte, val interface{}) byte {
	return 0
}

func (this *ZbnpiApp) ProInReadRspCmd(srcAddr uint16, hdr *ltl.FrameHdr, rdRspStatus []ltl.ReadRspStatus, val interface{}) error {
	itm, ok := val.(*elmodels.ItemInfos)
	if !ok {
		return errors.New("val assert elmodels.CacheqItem failed")
	}
	switch hdr.TrunkID {
	case ltl.TrunkID_GeneralBasic:
		gba := limp.BasicAttribute(int(hdr.NodeNo), rdRspStatus)
		if itm.IsLocal {
			devmodels.UpdateZbDeviceAndNode(itm.Sn, srcAddr, 1, gba.ProductIdentifier)
			if IsNetworkSteering() {
				elinkpsh.DeviceAnnce(gba.ProductIdentifier, itm.Sn, true)
			}
			return nil
		}
		out, err := jsoniter.Marshal(
			elmodels.DevPropRspPy{
				Sn:        itm.Sn,
				ProductID: itm.ProductID,
				Data:      gba,
			})
		if err != nil {
			logs.Debug(err)
			return err
		}
		ctrl.WriteResponse(itm.Client, itm.Tp, elink.CodeSuccess, itm.Pkid, out)
	default:
		return nil
	}

	return nil
}

func (this *ZbnpiApp) ProInWriteRspCmd(srcAddr uint16, hdr *ltl.FrameHdr, wrStatus []ltl.WriteRspStatus, val interface{}) error {
	return nil
}

func (this *ZbnpiApp) ProInWriteRpCfgRspCmd(srcAddr uint16, hdr *ltl.FrameHdr, crStatus []ltl.WriteRpCfgRspStatus, val interface{}) error {
	return nil
}

func (this *ZbnpiApp) ProInReadRpCfgRspCmd(srcAddr uint16, hdr *ltl.FrameHdr, rcStatus []ltl.ReadRpCfgRspStatus, val interface{}) error {
	return nil
}

func (this *ZbnpiApp) ProInReportCmd(srcAddr uint16, hdr *ltl.FrameHdr, rRec []ltl.ReportRec) error {
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

func (this *ZbnpiApp) ProInDefaultRsp(srcAddr uint16, hdr *ltl.FrameHdr, dfStatus *ltl.DefaultRsp, val interface{}) error {
	return nil
}
