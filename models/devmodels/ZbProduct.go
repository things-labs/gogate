package devmodels

import (
	"github.com/slzm40/gomo/ltl"
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
	ProductID_Switch: &ZbProduct{[]NodeDsc{
		{[]uint16{ltl.TrunkID_GeneralBasic, ltl.TrunkID_GeneralOnoff}, []uint16{}}}, "swtich"}, // 开关

	// 例子
	80000: &ZbProduct{[]NodeDsc{
		{[]uint16{}, []uint16{}},             // 节点1 集描述
		{[]uint16{1, 2}, []uint16{3, 4}},     // 节点2 集描述
		{[]uint16{3}, []uint16{5, 6, 7, 8}}}, //  节点3 集描述
		"测试1"},
}

// 根据产品id查找产品
func LookupZbProduct(pID int) (*ZbProduct, error) {
	o, exists := zbProduct[pID]
	if !exists {
		return nil, ErrProductNotExist
	}

	return o, nil
}

// 是否有这个产品Id
func HasZbProduct(pID int) bool {
	_, exists := zbProduct[pID]
	return exists
}

// 根据产品Id获得产品的节点描述,不含保留默认节点0
func LookupZbProductDeviceNodeDscList(pID int) ([]NodeDsc, error) {
	pdt, err := LookupZbProduct(pID)
	if err != nil {
		return nil, err
	}

	return pdt.GetDeviceNodeDscList(), nil
}

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
