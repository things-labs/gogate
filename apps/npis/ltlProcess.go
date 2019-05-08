package npis

import (
	"errors"
	"fmt"

	"github.com/thinkgos/gogate/apps/broad"
	"github.com/thinkgos/gogate/controllers/elinkpsh"
	"github.com/thinkgos/gogate/middle/ewait"
	"github.com/thinkgos/gogate/models"
	"github.com/thinkgos/gogate/protocol/elinkch/ctrl"
	"github.com/thinkgos/gogate/protocol/elinkmd"
	"github.com/thinkgos/gomo/elink"
	"github.com/thinkgos/gomo/ltl"
	"github.com/thinkgos/gomo/protocol/limp"

	"github.com/astaxie/beego/logs"
)

type deviceAnnce struct {
	sn      string
	nwkaddr uint16
}

func (this *ZbnpiApp) ProInSpecificCmd(srcAddr uint16, hdr *ltl.FrameHdr, cmdFormart []byte, val interface{}) byte {
	return 0
}

func (this *ZbnpiApp) ProInReadRspCmd(srcAddr uint16, hdr *ltl.FrameHdr, rdRspStatus []ltl.ReadRspStatus, val interface{}) error {
	var islocal bool
	id, ok := val.(string)
	if !ok {
		islocal = true
	}
	switch hdr.TrunkID {
	case ltl.TrunkID_GeneralBasic:
		gba := limp.BasicAttribute(rdRspStatus)
		if islocal {
			itm, ok := val.(deviceAnnce)
			if !ok || srcAddr != itm.nwkaddr {
				return errors.New("no this address")
			}
			logs.Debug("New device-productID: %d,sn: %s, srcAddress: %d",
				gba.ProductIdentifier, itm.sn, srcAddr)
			err := models.UpdateZbDeviceAndNode(itm.sn, srcAddr, 1, gba.ProductIdentifier)
			if err != nil {
				return err
			}
			if IsNetworkSteering() {
				elinkpsh.DeviceAnnce(gba.ProductIdentifier, itm.sn, true)
			}
			return nil
		}
		ewait.Done(id, gba)
	default:
		return errors.New(fmt.Sprintf("trunk not implementation: %d", hdr.TrunkID))
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

type ProReportPayload struct {
	ProductID int         `json:"productID"`
	Sn        string      `json:"sn"`
	NodeNo    int         `json:"nodeNo,omitempty"`
	Data      interface{} `json:"data"`
}

func (this *ZbnpiApp) ProInReportCmd(srcAddr uint16, hdr *ltl.FrameHdr, rRec []ltl.ReportRec) error {
	var v ProReportPayload

	zbdnode, err := models.LookupZbDeviceNodeByNN(srcAddr, hdr.NodeNo)
	if err != nil {
		return err
	}

	switch hdr.TrunkID {
	case ltl.TrunkID_GeneralOnoff:
		v = ProReportPayload{
			ProductID: zbdnode.ProductId,
			Sn:        zbdnode.Sn,
			NodeNo:    int(hdr.NodeNo),
			Data:      limp.OnoffAttribute(rRec),
		}
	case ltl.TrunkID_MsTemperatureMeasurement, ltl.TrunkID_MsRelativeHumidity:
		ms, err := limp.MsMeasureAttrReport(hdr.TrunkID, rRec)
		if err != nil {
			return err
		}

		v = ProReportPayload{
			ProductID: zbdnode.ProductId,
			Sn:        zbdnode.Sn,
			NodeNo:    int(hdr.NodeNo),
			Data:      ms,
		}
	default:
		logs.Error("no fix trunkid")
		return errors.New("no fix trunkid")
	}

	tp := elink.FormatPshTopic(ctrl.ChannelData,
		elink.FormatResouce(elinkmd.DevicePropertys, zbdnode.ProductId),
		elink.MethodPatch, elink.MessageTypeAnnce)

	return broad.PublishPyServerJSON(tp, v)
}

func (this *ZbnpiApp) ProInDefaultRsp(srcAddr uint16, hdr *ltl.FrameHdr, dfStatus *ltl.DefaultRsp, val interface{}) error {
	return nil
}
