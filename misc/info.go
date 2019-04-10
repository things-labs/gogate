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
	major  = 0
	minor  = 0
	fixed  = 1
	isBeta = true
)

type gatewayInfo struct {
	mac     string // "0C5415B171AA"
	MAC     string // "0C:54:15:B1:71:AA"
	version string // v1.2.3 Beta
}

var gatewayInfos *gatewayInfo

func init() {
	inf, err := NetInterface()
	if err != nil {
		panic(err)
	}

	gatewayInfos = &gatewayInfo{
		mac:     strings.ToUpper(hex.EncodeToString(inf.HardwareAddr)),
		MAC:     strings.ToUpper(inf.HardwareAddr.String()),
		version: utils.Version(major, minor, fixed, isBeta),
	}
}

func Mac() string {
	return gatewayInfos.mac
}

func MAC() string {
	return gatewayInfos.MAC
}
func Version() string {
	return gatewayInfos.version
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
