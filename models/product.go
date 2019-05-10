package models

// ProductInfo 产品信息
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

// RegisterProducts 注册产品列表
func RegisterProducts(pds map[int]*ProductInfo) {
	for k, v := range pds {
		_ = RegisterProduct(k, v)
	}
}

// RegisterProduct 注册相应产品
func RegisterProduct(pid int, pi *ProductInfo) error {
	if pid == PidRESERVE || pi == nil {
		return ErrInvalidParameter
	}
	productInfos[pid] = pi
	return nil
}

// LookupProduct 查找对应pid产品信息,types指定的产品类型
func LookupProduct(pid int, types ...int) (*ProductInfo, error) {
	if pid == PidRESERVE {
		return nil, ErrInvalidParameter
	}
	v, exist := productInfos[pid]
	if exist && (len(types) == 0 || (v.Types == types[0])) {
		return v, nil
	}

	return nil, ErrProductNotExist
}

// HasProduct 判断对应id产品信息是否存在,types指定的产品类型
func HasProduct(pid int, types ...int) bool {
	_, err := LookupProduct(pid, types...)
	return err == nil
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
