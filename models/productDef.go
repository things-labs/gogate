package models

import "github.com/thinkgos/gomo/ltl"

// 设备产品类型
const (
	PTypesGeneral = iota // 通用产品
	PTypesZigbee         // zigbee
	PTypesModbus         // modbus
)

// 所有的产品id列表, 必需注册到DeviceProductInfos, zigbee的产品需要另外注册到zigbee的设备产品里
const (
	PidDZMS01  = 20000 + iota // LC_DZMS01型号温湿度传感器
	PidDZSW01                 // LC_DZSW01型号一位智能开关
	PidDZSW02                 // LC_DZSW02型号二位智能开关
	PidDZSW03                 // LC_DZSW03型号三位智能开关
	PidDZCT01                 // LC_DZCT01型号窗帘控制器
	PidRESERVE = 0            // 保留pid号
)

// 产品表
var productInfos = map[int]*ProductInfo{
	PidDZMS01: &ProductInfo{0, PTypesZigbee, "smart zigbee",
		"LC_DZMS01", "温湿度传感器", "温湿度传感器", "lchtime",
		[]NodeDsc{
			{[]uint16{ltl.TrunkID_MsTemperatureMeasurement}, []uint16{}},
			{[]uint16{ltl.TrunkID_MsRelativeHumidity}, []uint16{}}}},

	PidDZSW01: &ProductInfo{0, PTypesZigbee, "smart zigbee",
		"LC_DZSW01", "一位智能开关", "一开智能开关", "lchtime",
		[]NodeDsc{
			{[]uint16{ltl.TrunkID_GeneralOnoff}, []uint16{}}}},

	PidDZSW02: &ProductInfo{0, PTypesZigbee, "smart zigbee",
		"LC_DZSW02", "二位智能开关", "二位智能开关", "lchtime",
		[]NodeDsc{
			{[]uint16{ltl.TrunkID_GeneralOnoff}, []uint16{}},
			{[]uint16{ltl.TrunkID_GeneralOnoff}, []uint16{}}}},

	PidDZSW03: &ProductInfo{0, PTypesZigbee, "smart zigbee",
		"LC_DZSW03", "三位智能开关", "三位智能开关", "lchtime",
		[]NodeDsc{
			{[]uint16{ltl.TrunkID_GeneralOnoff}, []uint16{}},
			{[]uint16{ltl.TrunkID_GeneralOnoff}, []uint16{}},
			{[]uint16{ltl.TrunkID_GeneralOnoff}, []uint16{}}}},

	PidDZCT01: &ProductInfo{0, PTypesZigbee, "smart zigbee",
		"LC_DZCT01", "窗帘控制器", "窗帘控制器", "lchtime",
		[]NodeDsc{
			{[]uint16{ltl.TrunkID_GeneralOnoff}, []uint16{}}}},
}
