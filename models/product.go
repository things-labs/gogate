package models

import (
    "github.com/jinzhu/gorm"
    "github.com/json-iterator/go"
)

type Product struct {
    gorm.Model
    ProductId uint32 `gorm:"UNIQUE;NOT NULL"`                                                         // 产品Id编号
    NodeList  string `gorm:'NOT NULL;default:"{\"NodeDscList\":[{\"InTrunk\":[],\"OutTrunk\":[]}]}"'` // 节点表
}

type NodeDsc struct {
    InTrunk  []uint16 // 输入集列表
    OutTrunk []uint16 //输出集列表
}

type NodeTables struct {
    NodeDscList []*NodeDsc // 节点描述列表
}

func FindProduct(pID uint32) (*Product, error) {
    o := &Product{}
    if devDb.Where(&Product{ProductId: pID}).First(o).RecordNotFound() == true {
        return nil, gorm.ErrRecordNotFound
    }

    return o, nil
}

func AddNewProduct(pID uint32, node []*NodeDsc) error {
    _, err := FindProduct(pID)
    if err != nil {
        newPdt := &Product{ProductId: pID}
        if node != nil {
            newPdt.setNodeDscList(node)
        }
        return devDb.Create(newPdt).Error
    }

    return nil
}
func DeleteProduct(pID uint32) error {
    o, err := FindProduct(pID)
    if err != nil {
        return nil
    }

    return devDb.Unscoped().Delete(o).Error

}

func GetProductNodeDscList(pID uint32) ([]*NodeDsc, error) {
    pdt, err := FindProduct(pID)
    if err != nil {
        return nil, err
    }

    return pdt.GetNodeDscList()
}

func (this *Product) GetNodeDscList() ([]*NodeDsc, error) {
    tb := &NodeTables{}
    if err := jsoniter.UnmarshalFromString(this.NodeList, tb); err != nil {
        return nil, err
    }

    return tb.NodeDscList, nil
}

func (this *Product) setNodeDscList(dsc []*NodeDsc) error {
    var err error

    tb := &NodeTables{NodeDscList: dsc}
    this.NodeList, err = jsoniter.MarshalToString(tb)
    if err != nil {
        return err
    }

    return nil
}

func (this *NodeDsc) SetTrunk(inTrunk, outTrunk []uint16) {
    this.InTrunk = append(this.InTrunk, inTrunk...)
    this.OutTrunk = append(this.OutTrunk, outTrunk...)
}

func (this *NodeDsc) GetTrunk() (inTrunk, outTrunk []uint16) {
    return this.InTrunk, this.OutTrunk
}
