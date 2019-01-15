package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/json-iterator/go"
)

type Product struct {
	gorm.Model
	ProductId   uint32 `gorm:"UNIQUE;NOT NULL"`                                                         // 产品Id编号
	NodeList    string `gorm:'NOT NULL;default:"{\"NodeDscList\":[{\"InTrunk\":[],\"OutTrunk\":[]}]}"'` // 节点表
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
		if node != nil {
			newPdt.SetNodeDscList(node)
		}
		return devDb.Create(newPdt).Error
	}

	return nil
}

// 增加一个新产品,提供所有节点描述
func (this *Product) AddProduct() error {
	_, err := LookupProduct(this.ProductId)
	if err != nil {
		return devDb.Create(this).Error
	}

	return nil
}

func UpdateProductDescritption(pID uint32, Newdesc string) error {
	o, err := LookupProduct(pID)
	if err != nil {
		return err
	}

	return devDb.Model(o).Update("description", Newdesc).Error
}

// 删除一个产品
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

// 获得产品的所有节点描述,不含保留默认节点0
func (this *Product) GetDeviceNodeDscList() ([]*NodeDsc, error) {
	if len(this.NodeList) == 0 {
		return nil, errors.New("no device node")
	}

	tb := &NodeTables{}
	if err := jsoniter.UnmarshalFromString(this.NodeList, tb); err != nil {
		return nil, err
	}

	return tb.NodeDscList, nil
}

//设置产品的所有节点描述,不含保留默认节点0
func (this *Product) SetNodeDscList(dsc []*NodeDsc) error {
	var err error

	tb := &NodeTables{NodeDscList: dsc}
	this.NodeList, err = jsoniter.MarshalToString(tb)
	if err != nil {
		return err
	}

	return nil
}

// 设置节点描述的输入输出集
func (this *NodeDsc) SetTrunk(inTrunk, outTrunk []uint16) {
	this.InTrunk = append(this.InTrunk, inTrunk...)
	this.OutTrunk = append(this.OutTrunk, outTrunk...)
}

// 获得节点描述的输入输出集
func (this *NodeDsc) GetTrunk() (inTrunk, outTrunk []uint16) {
	return this.InTrunk, this.OutTrunk
}
