package npis

import (
	"errors"
	"time"

	"github.com/thinkgos/gogate/misc"
	"github.com/thinkgos/gomo/ltl"
	"github.com/thinkgos/gomo/npi"
	"github.com/thinkgos/memlog"

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

func NewSerialConfig() *serial.Config {
	cfg := misc.APPConfig.Com0
	memlog.Debug(cfg)
	parity := serial.Parity('N')
	switch cfg.Parity {
	case "O":
		parity = serial.ParityOdd
	case "E":
		parity = serial.ParityEven
	case "M":
		parity = serial.ParityMark
	case "S":
		parity = serial.ParitySpace
	}

	return &serial.Config{
		Name:     cfg.Name,
		Baud:     cfg.BaudRate,
		Size:     byte(cfg.DataBit),
		Parity:   parity,
		StopBits: serial.StopBits(cfg.StopBit),
	}
}

func OpenZbApp() error {
	monitor, err := npi.Open(NewSerialConfig())
	if err != nil {
		return err
	}

	mid := NewMiddleMonitor(monitor)
	ZbApps := &ZbnpiApp{
		Ltl_t:         ltl.NewClient(mid),
		MiddleMonitor: mid,
	}

	go ZbApps.ServerInApdu(ZbApps.Context(), ZbApps)
	return ZbApps.NetworkFormation()
}

func CloseZbApp() {
	ZbApps.Close()
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
