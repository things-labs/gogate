package models

import "github.com/thinkgos/gomo/ltl"

// 设备产品类型
const (
	PTypesGeneral = iota // 通用产品
	PTypesZigbee         // zigbee产品
	PTypesModbus
)

// 所有的产品id列表, 必需注册到DeviceProductInfos, zigbee的产品需要另外注册到zigbee的设备产品里
const (
	PID_DZMS01  = 20000 + iota // LC_DZMS01型号温湿度传感器
	PID_DZSW01                 // LC_DZSW01型号一位智能开关
	PID_DZSW02                 // LC_DZSW02型号二位智能开关
	PID_DZSW03                 // LC_DZSW03型号三位智能开关
	PID_DZCT01                 // LC_DZCT01型号窗帘控制器
	PID_RESERVE = 0            // 保留pid号
)

// 产品信息
type ProductInfo struct {
	Number           int         // 编号  用于识别相同的类型走的通道 0为默认通道
	Types            int         // 类型
	TypesName        string      // 类型名称
	ModelSpec        string      // 型号规格
	ModelName        string      // 型号名称
	Description      string      // 描述
	ManufacturerName string      // 制造商名字
	Attach           interface{} // 附加信息, 用于定义产品属性参数等
}

var productInfos = map[int]*ProductInfo{
	PID_DZMS01: &ProductInfo{0, PTypesZigbee, "smart zigbee",
		"LC_DZMS01", "温湿度传感器", "温湿度传感器", "lchtime",
		[]NodeDsc{
			{[]uint16{ltl.TrunkID_MsTemperatureMeasurement}, []uint16{}},
			{[]uint16{ltl.TrunkID_MsRelativeHumidity}, []uint16{}}}},

	PID_DZSW01: &ProductInfo{0, PTypesZigbee, "smart zigbee",
		"LC_DZSW01", "一位智能开关", "一开智能开关", "lchtime",
		[]NodeDsc{
			{[]uint16{ltl.TrunkID_GeneralOnoff}, []uint16{}}}},

	PID_DZSW02: &ProductInfo{0, PTypesZigbee, "smart zigbee",
		"LC_DZSW02", "二位智能开关", "二位智能开关", "lchtime",
		[]NodeDsc{
			{[]uint16{ltl.TrunkID_GeneralOnoff}, []uint16{}},
			{[]uint16{ltl.TrunkID_GeneralOnoff}, []uint16{}}}},

	PID_DZSW03: &ProductInfo{0, PTypesZigbee, "smart zigbee",
		"LC_DZSW03", "三位智能开关", "三位智能开关", "lchtime",
		[]NodeDsc{
			{[]uint16{ltl.TrunkID_GeneralOnoff}, []uint16{}},
			{[]uint16{ltl.TrunkID_GeneralOnoff}, []uint16{}},
			{[]uint16{ltl.TrunkID_GeneralOnoff}, []uint16{}}}},

	PID_DZCT01: &ProductInfo{0, PTypesZigbee, "smart zigbee",
		"LC_DZCT01", "窗帘控制器", "窗帘控制器", "lchtime",
		[]NodeDsc{
			{[]uint16{ltl.TrunkID_GeneralOnoff}, []uint16{}}}},
}

// RegisterProducts 注册产品列表
func RegisterProducts(pds map[int]*ProductInfo) {
	for k, v := range pds {
		_ = RegisterProduct(k, v)
	}
}

// RegisterProduct 注册相应产品
func RegisterProduct(pid int, pi *ProductInfo) error {
	if pid == PID_RESERVE || pi == nil {
		return ErrInvalidParameter
	}
	productInfos[pid] = pi
	return nil
}

// LookupProduct 查找对应pid产品信息
func LookupProduct(pid int) (*ProductInfo, error) {
	if pid == PID_RESERVE {
		return nil, ErrInvalidParameter
	}
	if v, exist := productInfos[pid]; exist {
		return v, nil
	}
	return nil, ErrProductNotExist
}

// HasProduct 判断对应id产品信息是否存在,types指定的产品类型
func HasProduct(pid int, types ...int) bool {
	info, err := LookupProduct(pid)
	if err != nil {
		return false
	}
	if len(types) > 0 {
		return info.Types == types[0]
	}
	return true
}

/*****************************************************************************
				ZIGBEE ONLY
******************************************************************************/

// NodeDsc 节点输入输出集列表
type NodeDsc struct {
	InTrunk  []uint16
	OutTrunk []uint16
}

// 获取产品的所有节点描述,不含保留默认节点0
func (this *ProductInfo) GetZbDeviceNodeDscList() []NodeDsc {
	return this.Attach.([]NodeDsc)
}
