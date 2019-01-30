package models

import (
	"errors"
)

type NodeDsc struct {
	InTrunk  []uint16 // 输入集列表
	OutTrunk []uint16 // 输出集列表
}

type Product struct {
	NodeList    []NodeDsc
	Description string
}

var product map[uint32]*Product = map[uint32]*Product{
	// ProductID: 节点列表,节点描述
	3000: &Product{NodeList: []NodeDsc{{InTrunk: []uint16{}, OutTrunk: []uint16{}}, {InTrunk: []uint16{1, 2}, OutTrunk: []uint16{3, 4}}, {InTrunk: []uint16{3}, OutTrunk: []uint16{5, 6, 7, 8}}}, Description: "测试1"},
}

// 根据产品id查找产品
func LookupProduct(pID uint32) (*Product, error) {

	o, exists := product[pID]
	if !exists {
		return nil, errors.New("product no exist")
	}

	return o, nil
}

// 根据产品Id获得产品的所有节点描述,不含保留默认节点0
func LookupProductDeviceNodeDscList(pID uint32) ([]NodeDsc, error) {
	pdt, err := LookupProduct(pID)
	if err != nil {
		return nil, err
	}

	return pdt.GetDeviceNodeDscList()
}

func (this *Product) GetProductDescritption() string {
	return this.Description
}

// 获取产品的所有节点描述,不含保留默认节点0
func (this *Product) GetDeviceNodeDscList() ([]NodeDsc, error) {
	return this.NodeList, nil
}

// 获取节点描述的输入输出集
func (this *NodeDsc) GetTrunkID() (inTrunk, outTrunk []uint16) {
	return this.InTrunk, this.OutTrunk
}

/*

import (
	"github.com/jinzhu/gorm"
	"github.com/json-iterator/go"
)

const _default_product_node_list = `{"NodeDscList":[]}`

type Product struct {
	gorm.Model
	ProductId   uint32 `gorm:"UNIQUE;NOT NULL"`   // 产品Id编号
	NodeList    string `gorm:"type:varchar(511)"` // 节点表
	Description string
}

type NodeDsc struct {
	InTrunk  []uint16 // 输入集列表
	OutTrunk []uint16 // 输出集列表
}

type NodeTables struct {
	NodeDscList []*NodeDsc // 节点描述列表
}

// 根据产品id查找产品
func LookupProduct(pID uint32) (*Product, error) {
	o := &Product{}
	if devDb.Where(&Product{ProductId: pID}).First(o).RecordNotFound() == true {
		return nil, gorm.ErrRecordNotFound
	}

	return o, nil
}

// 增加一个新产品,提供所有节点描述
func AddProduct(pID uint32, node []*NodeDsc, desc string) error {
	_, err := LookupProduct(pID)
	if err != nil {
		newPdt := &Product{ProductId: pID, Description: desc}
		newPdt.SetNodeDscList(node)
		return devDb.Create(newPdt).Error
	}

	return nil
}

// 增加一个新产品,提供所有节点描述
func (this *Product) AddProduct() error {
	_, err := LookupProduct(this.ProductId)
	if err != nil {
		if len(this.NodeList) == 0 {
			this.NodeList = _default_product_node_list
		}
		return devDb.Create(this).Error
	}

	return nil
}

// 更新产品描述
func UpdateProductDescritption(pID uint32, Newdesc string) error {
	o, err := LookupProduct(pID)
	if err != nil {
		return err
	}

	return devDb.Model(o).Update("description", Newdesc).Error
}

//获取产品描述
func GetProductDescription(pID uint32) (string, error) {
	o, err := LookupProduct(pID)
	if err != nil {
		return "", err
	}

	return o.Description, nil
}

// 根据产品Id删除一个产品
func DeleteProduct(pID uint32) error {
	o, err := LookupProduct(pID)
	if err != nil {
		return nil
	}

	return devDb.Unscoped().Delete(o).Error
}

// 根据产品Id获得产品的所有节点描述,不含保留默认节点0
func LookupProductDeviceNodeDscList(pID uint32) ([]*NodeDsc, error) {
	pdt, err := LookupProduct(pID)
	if err != nil {
		return nil, err
	}

	return pdt.GetDeviceNodeDscList()
}

// 获取产品的所有节点描述,不含保留默认节点0
func (this *Product) GetDeviceNodeDscList() ([]*NodeDsc, error) {
	tb := &NodeTables{}
	if err := jsoniter.UnmarshalFromString(this.NodeList, tb); err != nil {
		return nil, err
	}

	return tb.NodeDscList, nil
}

//设置产品的所有节点描述,不含保留默认节点0
func (this *Product) SetNodeDscList(dsc []*NodeDsc) error {
	var err error

	tb := &NodeTables{}
	if dsc == nil {
		tb.NodeDscList = []*NodeDsc{}
	} else {
		tb.NodeDscList = dsc
	}

	if this.NodeList, err = jsoniter.MarshalToString(tb); err != nil {
		return err
	}

	return nil
}

// 设置节点描述的输入输出集
func (this *NodeDsc) SetTrunkID(inTrunk, outTrunk []uint16) {
	if inTrunk == nil {
		this.InTrunk = []uint16{}
	} else {
		this.InTrunk = inTrunk
	}

	if outTrunk == nil {
		this.OutTrunk = []uint16{}
	} else {
		this.OutTrunk = outTrunk
	}
}

// 获取节点描述的输入输出集
func (this *NodeDsc) GetTrunkID() (inTrunk, outTrunk []uint16) {
	return this.InTrunk, this.OutTrunk
}
*/
