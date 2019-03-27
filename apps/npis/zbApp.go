package npis

import (
	"errors"
	"time"

	"github.com/slzm40/gomo/ltl"
	"github.com/slzm40/gomo/misc"
	"github.com/slzm40/gomo/npi"

	"github.com/astaxie/beego/logs"
	"github.com/tarm/serial"
)

const Incomming_msg_size_max = 256

type ZbnpiApp struct {
	isNetworkFormation bool
	isNetworkSteering  bool
	*ltl.Ltl_t
	*MiddleMonitor
}

var ZbApps *ZbnpiApp

func ZbAppInit() error {
	var err error
	var m *npi.Monitor

	bcfg := misc.APPCfg
	usartcfg := &serial.Config{}

	if usartcfg.Name, err = bcfg.GetValue("COM0", "Name"); err != nil {
		return err
	}

	usartcfg.Baud = bcfg.MustInt("COM0", "Name", 115200)
	usartcfg.Size = byte(bcfg.MustInt("COM0", "DataBit", 8))
	usartcfg.Parity = serial.Parity(bcfg.MustInt("COM0", "Parity", 'N'))
	usartcfg.StopBits = serial.StopBits(bcfg.MustInt("COM0", "StopBit", 1))

	logs.Debug("usarcfg: %#v", usartcfg)

	if m, err = npi.NewNpiMonitor(usartcfg); err != nil {
		logs.Error("npi new failed", err)
		return err
	}

	mid := &MiddleMonitor{
		IncommingMsgPkt: make(chan *ltl.IncomingMsgPkt, Incomming_msg_size_max),
		Monitor:         m,
	}

	mid.AddAsyncCbs(map[uint16]func(*npi.Npdu){
		npi.MT_AF_DATA_CONFIRM:                        Af_DataConfirm,
		npi.MT_AF_INCOMING_MSG:                        Af_IncomingMsg,
		npi.MT_ZDO_MGMT_PERMIT_JOIN_RSP:               Zdo_MgmtPermitJoinRsp,
		npi.MT_ZDO_STATE_CHANGE_IND:                   Zdo_StateChangeInd,
		npi.MT_ZDO_END_DEVICE_ANNCE_IND:               Zdo_EnddeviceAnnceInd,
		npi.MT_ZDO_LEAVE_IND:                          Zdo_LeaveInd,
		npi.MT_SYS_RESET_IND:                          Sys_ResetInd,
		npi.MT_APP_CNF_BDB_COMMISSIONING_NOTIFICATION: Appcfg_BdbCommissioningNotice,
	})

	ZbApps = &ZbnpiApp{
		Ltl_t:         &ltl.Ltl_t{mid},
		MiddleMonitor: mid,
	}

	go ZbApps.ServerInApdu(ZbApps)
	ZbApps.MiddleMonitor.Start()

	return ZbApps.NetworkFormation()
}

// 建立zigbee的网络
func (this *ZbnpiApp) NetworkFormation() error {
	for trycnt := 0; ; trycnt++ {
		if ok, err := this.Appcfg_BdbStartCommissioningReq(
			npi.Cms_mode_NetworkFormation); err != nil || !ok {
			if trycnt == 10 {
				return errors.New("npis: Formation network failed")
			}
			time.Sleep(time.Millisecond * 500)
			continue
		} else {
			break
		}
	}

	return nil
}

func IsNetworkFormation() bool {
	return ZbApps.isNetworkFormation
}

func SetNetworkSteering(on bool) {
	ZbApps.isNetworkSteering = on
}

func IsNetworkSteering() bool {
	return ZbApps.isNetworkSteering
}
