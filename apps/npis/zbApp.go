package npis

import (
	"errors"
	"time"

	"github.com/thinkgos/gogate/misc"
	"github.com/thinkgos/gomo/ltl"
	"github.com/thinkgos/gomo/npi"

	"github.com/astaxie/beego/logs"
	"github.com/tarm/serial"
	"go.uber.org/dig"
)

const Incomming_msg_size_max = 256

type ZbnpiApp struct {
	isNetworkFormation bool
	isNetworkSteering  bool
	*ltl.Ltl_t
	*MiddleMonitor
}

var ZbApps *ZbnpiApp

func NewSerialConfig() (*serial.Config, error) {
	bcfg := misc.UartCfg
	usartcfg := &serial.Config{}

	secCom0, err := bcfg.GetSection("COM0")
	if err != nil {
		return nil, err
	}

	usartcfg.Name = secCom0.Key("Name").MustString("COM0")
	usartcfg.Baud = secCom0.Key("BaudRate").MustInt(115200)
	usartcfg.Size = byte(secCom0.Key("DataBit").MustUint(8))
	usartcfg.Parity = serial.Parity(secCom0.Key("Parity").MustInt('N'))
	usartcfg.StopBits = serial.StopBits(secCom0.Key("StopBit").MustInt(1))
	logs.Debug("usarcfg: %#v", usartcfg)

	return usartcfg, nil
}

func OpenZbApp() error {
	container := dig.New()
	container.Provide(NewSerialConfig)
	container.Provide(npi.Open)
	container.Provide(NewMiddleMonitor)
	container.Provide(func(mid *MiddleMonitor) *ZbnpiApp {
		return &ZbnpiApp{
			Ltl_t:         ltl.NewClient(mid),
			MiddleMonitor: mid,
		}
	})

	return container.Invoke(func(app *ZbnpiApp) {
		go app.ServerInApdu(app.Context(), app)
		app.NetworkFormation()
		ZbApps = app
	})
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
