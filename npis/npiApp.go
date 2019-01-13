package npis

import (
	"github.com/Unknwon/goconfig"
	"github.com/astaxie/beego/logs"
	"github.com/slzm40/gomo/npi"
	"github.com/tarm/serial"
)

const UART_CFG_PATH = "./conf/usart.conf" // 串口配置的路径
var nobj *npi.NpiObj

func DoneCmd(cmd uint16) {
	var err error

	logs.Info("0x%04x", cmd)

	pdu := &npi.Npi_pdu{}
	switch cmd {
	case npi.MT_SYS_PING:
		pdu.Sys_PingReq_Pack()
		if rsp, err := nobj.SendSynchData(pdu); err != nil {
			logs.Error(err)
		} else {
			if Capabilities, err := rsp.Sys_PingSRspParse(); err != nil {
				logs.Error(err)
			} else {
				logs.Debug("0x%04x", Capabilities)
			}
		}
	case npi.MT_SYS_RESET_REQ:
		pdu.Sys_ResetReq_Pack(npi.MT_SYS_RESET_HARD)
		if err = nobj.SendAsynchData(pdu); err != nil {
			logs.Error(err)
		}
	case npi.MT_APP_CNF_BDB_START_COMMISSIONING:
		pdu.Appcfg_BdbStartCommissioningReq_Pack(0x04)
		//		if err = nobj.
	}
}
func NpiAppInit() error {
	var bcfg *goconfig.ConfigFile
	var err error

	if bcfg, err = goconfig.LoadConfigFile(UART_CFG_PATH); err != nil {
		logs.Error("Npi Init: Load usart config file failed!")
		return err
	}

	usartcfg := &serial.Config{}
	if usartcfg.Name, err = bcfg.GetValue("COM0", "Name"); err != nil {
		return err
	}
	usartcfg.Baud = bcfg.MustInt("COM0", "Name", 115200)
	usartcfg.Size = byte(bcfg.MustInt("COM0", "DataBit", 8))
	usartcfg.Parity = serial.Parity(bcfg.MustInt("COM0", "Parity", 'N'))
	usartcfg.StopBits = serial.StopBits(bcfg.MustInt("COM0", "StopBit", 1))

	logs.Debug("usarcfg: %#v", usartcfg)

	if nobj, err = npi.NpiNew(usartcfg); err != nil {
		logs.Error("npi new failed", err)
		return err
	}

	return nil
}
