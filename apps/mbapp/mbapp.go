package mbapp

import (
	"time"

	"github.com/goburrow/serial"
	"github.com/thinkgos/gogate/middle/mb"
	modbus "github.com/thinkgos/gomodbus"
)

func init() {
	mbappInit()
}

var gatherPara1 = &mb.GatherPara{
	SlaveID:            0x01,
	HasCoil:            true,
	CoilAddress:        0,
	CoilQuantity:       2000,
	CoilVirtualAddress: 0,
	CoilScanRate:       1 * time.Second,

	HasDiscrete:            true,
	DiscreteAddress:        0,
	DiscreteQuantity:       2000,
	DiscreteVirtualAddress: 0,
	DiscreteScanRate:       1 * time.Second,

	HasInput:            true,
	InputAddress:        0,
	InputQuantity:       200,
	InputVirtualAddress: 0,
	InputScanRate:       1 * time.Second,

	HasHolding:            true,
	HoldingAddress:        0,
	HoldingQuantity:       200,
	HoldingVirtualAddress: 0,
	HoldingScanRate:       1 * time.Second,
}

func mbappInit() error {
	p := modbus.NewRTUClientProvider("COM1")
	p.Config = serial.Config{
		Address:  "COM1",
		BaudRate: 115200,
		DataBits: 8,
		StopBits: 1,
		Parity:   "N",
		Timeout:  time.Second * 1,
	}
	port1 := mb.NewClient(p)

	if err := port1.AddReadPoll(gatherPara1); err != nil {
		return err
	}
	err := port1.Start()

	srv := modbus.NewTCPServer(":502")
	srv.AddNode(port1.GetNodeRegister())
	go srv.ServerModbus()
	return err
}
