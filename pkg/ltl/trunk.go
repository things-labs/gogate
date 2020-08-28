package ltl

import (
	"errors"
	"strings"
)

const (
	NodeNumReserved = 0
)

// 基本集
const (
	TrunkID_GeneralBasic = iota
	TrunkID_GeneralPower
	TrunkID_GeneralOnoff
	TrunkID_GeneralLevelControl
)

// Measurement and Sensing trunks  测量与传感器集
const (
	TrunkID_MsIlluminanceMeasurement        = 0x0200 // 光照测量
	TrunkID_MsIlluminanceLevelSensingConfig = 0x0201 // 光照水平感知配置
	TrunkID_MsTemperatureMeasurement        = 0x0202 // 温度测量
	TrunkID_MsPressureMeasurement           = 0x0203 // 压力测量
	TrunkID_MsFlowMeasurement               = 0x0204 // 流量测量
	TrunkID_MsRelativeHumidity              = 0x0205 // 相对湿度测量
	TrunkID_MsOccupancySensing              = 0x0206 // 占有率测量
)

var Attribute map[uint16][]string = map[uint16][]string{
	0: {"LTLVersion", "APPVersion", "HWVersion", "ManufacturerName",
		"BuildDateCode", "ProductIdentifier", "SerialNumber", "PowerSource"}, // basic truck
	1: {},               // power trunk
	2: {"onoff"},        // onoff trunk
	3: {"currentLevel"}, // level control trunk
}

// 获得此集的属性id顶值
func AttrIDMax(trunkID uint16) (uint16, error) {
	if v, ok := Attribute[trunkID]; ok {
		return uint16(len(v)), nil
	}

	return 0, errors.New("not find the trunk")
}

// 通过属性name获得属性Id
func AttrID(trunkID uint16, Attrstr string) (uint16, error) {
	if v, ok := Attribute[trunkID]; ok {
		for i, vid := range v {
			if strings.EqualFold(Attrstr, vid) {
				return uint16(i), nil
			}
		}
	}

	return 0, errors.New("not find the attribute id")
}

// 通用属性id获得属性Name
func AttrStr(trunkID uint16, AttrId uint16) (string, error) {
	if v, ok := Attribute[trunkID]; ok && int(AttrId) < len(v) {
		return v[AttrId], nil
	}

	return "", errors.New("not find the attribute id")
}
