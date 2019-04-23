package models

import (
	"github.com/thinkgos/gomo/ltl"
)

// 节点输入输出集列表
type NodeDsc struct {
	InTrunk  []uint16
	OutTrunk []uint16
}

// zigbee 产品节点描述
type ZbProduct struct {
	NodeList    []NodeDsc
	Description string
}

var zbProduct map[int]*ZbProduct = map[int]*ZbProduct{
	// ProductID: 节点列表,节点描述
	PID_DZMS01: &ZbProduct{[]NodeDsc{
		{[]uint16{ltl.TrunkID_MsTemperatureMeasurement}, []uint16{}},
		{[]uint16{ltl.TrunkID_MsRelativeHumidity}, []uint16{}}},
		"温湿度传感器"},
	PID_DZSW01: &ZbProduct{[]NodeDsc{
		{[]uint16{ltl.TrunkID_GeneralOnoff}, []uint16{}}}, //
		"一位智能开关"},
	PID_DZSW02: &ZbProduct{[]NodeDsc{
		{[]uint16{ltl.TrunkID_GeneralOnoff}, []uint16{}},
		{[]uint16{ltl.TrunkID_GeneralOnoff}, []uint16{}}},
		"二位智能开关"},
	PID_DZSW03: &ZbProduct{[]NodeDsc{
		{[]uint16{ltl.TrunkID_GeneralOnoff}, []uint16{}},
		{[]uint16{ltl.TrunkID_GeneralOnoff}, []uint16{}},
		{[]uint16{ltl.TrunkID_GeneralOnoff}, []uint16{}}},
		"三位智能开关"},

	// 例子
	PID_ZIGBEE_TEST: &ZbProduct{[]NodeDsc{
		{[]uint16{}, []uint16{}},             // 节点1 集描述
		{[]uint16{1, 2}, []uint16{3, 4}},     // 节点2 集描述
		{[]uint16{3}, []uint16{5, 6, 7, 8}}}, //  节点3 集描述
		"测试1"},
}

// 根据产品id查找产品
func LookupZbProduct(pid int) (*ZbProduct, error) {
	if pid == PID_RESERVE {
		return nil, ErrInvalidParameter
	}
	if o, exists := zbProduct[pid]; exists {
		return o, nil
	}
	return nil, ErrProductNotExist
}

// 是否有这个产品Id
func HasZbProduct(pid int) bool {
	if pid == PID_RESERVE {
		return false
	}
	_, exists := zbProduct[pid]
	return exists
}

// 根据产品Id获得产品的节点描述,不含保留默认节点0
func GetDeviceNodeDscList(pid int) ([]NodeDsc, error) {
	pdt, err := LookupZbProduct(pid)
	if err != nil {
		return nil, err
	}

	return pdt.GetDeviceNodeDscList(), nil
}

// 获取产品描述
func (this *ZbProduct) GetProductDescritption() string {
	return this.Description
}

// 获取产品的所有节点描述,不含保留默认节点0
func (this *ZbProduct) GetDeviceNodeDscList() []NodeDsc {
	return this.NodeList
}

// 获取节点描述的输入输出集
func (this *NodeDsc) GetTrunkID() (inTrunk, outTrunk []uint16) {
	return this.InTrunk, this.OutTrunk
}
