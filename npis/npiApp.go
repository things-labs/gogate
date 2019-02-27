package npis

import (
	"github.com/astaxie/beego/logs"
	"github.com/slzm40/gomo/misc"
	"github.com/slzm40/gomo/npi"
	"github.com/tarm/serial"
)

const (
	zb_state_idle = iota
	zb_state_nwkFormation
)

type ZigbeeMonitor struct {
	*npi.Monitor
	state int
}

var ZbMonitor *ZigbeeMonitor

func DoneCmd(cmd uint16) {
	var err error

	logs.Info("cmd request: 0x%04x", cmd)

	switch cmd {
	case npi.MT_SYS_PING:
		if Capabilities, err := ZbMonitor.Sys_Ping(); err != nil {
			logs.Error(err)
		} else {
			logs.Debug("Capabilities: 0x%04x", Capabilities)
		}

	case npi.MT_SYS_RESET_REQ:
		if err = ZbMonitor.Sys_ResetReq(npi.MT_SYS_RESET_HARD); err != nil {
			logs.Error(err)
		}

	case npi.MT_APP_CNF_BDB_START_COMMISSIONING:
		if status, err := ZbMonitor.Appcfg_BdbStartCommissioningReq(0x04); err != nil {
			logs.Error(err)
		} else {
			logs.Debug("Appcfg_BdbStartCommissioningReq %t", status)
		}
	case 100:
		ZbMonitor.Close()
	}
}

func NpiAppInit() error {
	var (
		err error
		m   *npi.Monitor
	)

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

	ZbMonitor = &ZigbeeMonitor{
		Monitor: m,
		state:   zb_state_idle,
	}

	ZbMonitor.AddAsyncCbs(map[uint16]func(*npi.Npdu){
		npi.MT_AF_DATA_CONFIRM:                        Af_DataConfirm,
		npi.MT_AF_INCOMING_MSG:                        Af_IncomingMsgParse,
		npi.MT_ZDO_MGMT_PERMIT_JOIN_RSP:               Zdo_MgmtPermitJoinRsp,
		npi.MT_ZDO_STATE_CHANGE_IND:                   Zdo_StateChangeInd,
		npi.MT_ZDO_END_DEV_ANNCE:                      Zdo_enddeviceAnnceInd,
		npi.MT_ZDO_LEAVE_IND:                          Zdo_LeaveInd,
		npi.MT_SYS_RESET_IND:                          Sys_ResetInd,
		npi.MT_APP_CNF_BDB_COMMISSIONING_NOTIFICATION: Appcfg_BdbCommissioningNotice,
	})

	return nil
}

func ZigbeeNetworkFromationBoot() {

}
