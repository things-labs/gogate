package misc

import (
	"encoding/hex"
	"errors"
	"net"
	"runtime"
	"strings"

	"github.com/thinkgos/utils"
)

const (
	version = "v1.2.3 Beta"
)

type gatewayInfo struct {
	mac       string // "0C5415B171AA"
	MAC       string // "0C:54:15:B1:71:AA"
	version   string // v1.2.3 Beta
	buildDate string // 2018-12-09 15:26:26
}

var gatewayInfos *gatewayInfo

func init() {
	inf, err := NetInterface()
	if err != nil {
		panic(err)
	}

	gatewayInfos = &gatewayInfo{
		mac:       strings.ToUpper(hex.EncodeToString(inf.HardwareAddr)),
		MAC:       strings.ToUpper(inf.HardwareAddr.String()),
		version:   version,
		buildDate: utils.BuildDateTime(),
	}
}

// "0C5415B171AA"
func Mac() string {
	return gatewayInfos.mac
}

// "0C:54:15:B1:71:AA"
func MAC() string {
	return gatewayInfos.MAC
}

// v1.2.3 Beta
func Version() string {
	return gatewayInfos.version
}

// format 2018-12-09 15:26:26
func BuildDate() string {
	return gatewayInfos.buildDate
}

func NetInterface() (*net.Interface, error) {
	var intf *net.Interface
	var err error

	switch runtime.GOOS {
	case "windows":
		intf, err = net.InterfaceByName("WLAN") // 获取windows地址
		if err != nil {
			return nil, err
		}

	case "linux":
		intf, err = net.InterfaceByName("ens33") // 获取linux地址
		if err == nil {
			return intf, nil
		}

		intf, err = net.InterfaceByName("eth0")
		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("no os")
	}

	return intf, nil
}
