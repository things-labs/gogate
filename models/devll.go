package models

import (
	"os"
	"path"

	"github.com/Unknwon/com"
	"github.com/jinzhu/gorm"
	"github.com/json-iterator/go"
	_ "github.com/mattn/go-sqlite3"
)

const (
	_DB_NAME  = "data/devll.db"
	_DB_DRIVE = "sqlite3"
)

//设备表
type DeviceInfo struct {
	gorm.Model
	IeeeAddr  uint64 `gorm:"UNIQUE;NOT NULL"`
	NwkAddr   uint16 `gorm:"NOT NULL"`
	Capacity  byte   `gorm:"default:2"`
	ProductId uint32
}

//节点表
type NodeInfo struct {
	ID           uint   `gorm:"primary_key"`
	NwkAddr      uint16 `gorm:"UNIQUE;NOT NULL"`
	NodeNo       uint16 `gorm:"UNIQUE;NOT NULL"`
	IeeeAddr     uint64 `gorm:"NOT NULL"`
	InTrunkList  string `gorm:'default:"{\"trunkID\":[]}"'`
	OutTrunkList string `gorm:'default:"{\"trunkID\":[]}"'`
	BindList     string `gorm:'default:"{\"id\":[]}'`
}

type BindInfo struct {
	Id []uint `json:"id"`
}
type TrunkIDList struct {
	TrunkID []uint16 `json:"trunkID"`
}

var devDb *gorm.DB

func init() {
	var err error

	// 判断目录是否存在,不存在着创建对应的所有目录
	if !com.IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}

	if devDb, err = gorm.Open(_DB_DRIVE, _DB_NAME); err != nil {
		panic("models: gorm open failed," + err.Error())
	}
	//default disable
	//devDb.LogMode(misc.APPCfg.MustBool(goconfig.DEFAULT_SECTION, "ormDbLog", false))
	devDb.LogMode(true)

	devDb.AutoMigrate(&Product{}, &DeviceInfo{}, &NodeInfo{})
}

func FindDevllByNwk(nwkAddr uint16) (*DeviceInfo, error) {
	oInfo := &DeviceInfo{}
	if db := devDb.Where(&DeviceInfo{NwkAddr: nwkAddr}).First(oInfo); db.RecordNotFound() == true {
		return nil, gorm.ErrRecordNotFound
	}

	return oInfo, nil
}

func FindDevllByIeeeAddr(ieeeAddr uint64) (*DeviceInfo, error) {
	oInfo := &DeviceInfo{}
	if db := devDb.Where(&DeviceInfo{IeeeAddr: ieeeAddr}).First(oInfo); db.RecordNotFound() == true {
		return nil, gorm.ErrRecordNotFound
	}

	return oInfo, nil
}

func InserUpdateDevll(ieeeAddr uint64, nwkAddr uint16, capacity byte, productID uint32) error {
	oInfo, err := FindDevllByIeeeAddr(ieeeAddr)
	if err != nil {
		return devDb.Create(&DeviceInfo{
			IeeeAddr:  ieeeAddr,
			NwkAddr:   nwkAddr,
			Capacity:  capacity,
			ProductId: productID,
		}).Error
	}

	if (oInfo.NwkAddr == nwkAddr) && (oInfo.Capacity == capacity) && oInfo.ProductId == productID {
		return nil
	}

	oInfo.NwkAddr = nwkAddr
	oInfo.Capacity = capacity
	oInfo.ProductId = productID

	return devDb.Save(oInfo).Error
}

func DeleteDevll(ieeeAddr uint64) error {
	oInfo, err := FindDevllByIeeeAddr(ieeeAddr)
	if err != nil {
		return nil
	}

	return devDb.Unscoped().Delete(oInfo).Error
	//	return devDb.Delete(oInfo).Error
}

func (this *DeviceInfo) GetIeeeAddr() uint64 {
	return this.IeeeAddr
}
func (this *DeviceInfo) GetNwkAddr() uint16 {
	return this.NwkAddr
}
func (this *DeviceInfo) GetCapacity() byte {
	return this.Capacity
}
func (this *DeviceInfo) GetID() uint {
	return this.ID
}
func (this *DeviceInfo) GetProductID() uint32 {
	return this.ProductId
}

func FindNBI(nwkAddr, nodeNum uint16) (*NodeInfo, error) {
	oNbi := &NodeInfo{}
	if devDb.Where(&NodeInfo{NwkAddr: nwkAddr, NodeNo: nodeNum}).First(oNbi).RecordNotFound() == true {
		return nil, gorm.ErrRecordNotFound
	}

	return oNbi, nil
}

func FindNbiByID(id uint) (*NodeInfo, error) {
	oNbi := &NodeInfo{}
	if devDb.First(oNbi, id).RecordNotFound() == true {
		return nil, gorm.ErrRecordNotFound
	}

	return oNbi, nil
}

func InserUpdateNBI(nwkAddr, nodeNum uint16, ieeeAddr uint64, inTrunk, outTrunk string) error {
	oNbi, err := FindNBI(nwkAddr, nodeNum)
	if err != nil {
		return devDb.Create(&NodeInfo{
			NwkAddr:      nwkAddr,
			NodeNo:       nodeNum,
			InTrunkList:  inTrunk,
			OutTrunkList: outTrunk,
		}).Error
	}

	if len(inTrunk) != 0 {
		oNbi.InTrunkList = inTrunk
	}

	if len(outTrunk) != 0 {
		oNbi.OutTrunkList = outTrunk
	}

	return devDb.Save(oNbi).Error
}

func UpdateNBIBindList(nwkAddr, nodeNum, bindAddr, bindNodeNum uint16) error {
	//	var (
	//		srcNbi, bindNbi *NodeBindInfo
	//		err             error
	//	)
	//	if srcNbi, err = FindNBI(nwkAddr, nodeNum); err != nil {
	//		return errors.New("bind failed")
	//	}

	//	if bindNbi, err = FindNBI(bindAddr, bindNodeNum); err != nil {
	//		return errors.New("bind failed")
	//	}

	return nil
}

func (this *NodeInfo) GetNbiNwkAddr() uint16 {
	return this.NwkAddr
}
func (this *NodeInfo) GetNbiNodeNum() uint16 {
	return this.NodeNo
}
func (this *NodeInfo) GetNbiID() uint {
	return this.ID
}
func (this *NodeInfo) GetNbiTrunkIDList(isInTrunk bool) ([]uint16, error) {
	var str string

	if isInTrunk {
		str = this.InTrunkList
	} else {
		str = this.OutTrunkList
	}

	tmp := &TrunkIDList{}
	if err := jsoniter.UnmarshalFromString(str, tmp); err != nil {
		return nil, err
	}

	return tmp.TrunkID, nil
}

func (this *NodeInfo) GetNbiBindList() ([]uint, error) {
	o := &BindInfo{}
	err := jsoniter.UnmarshalFromString(this.BindList, o)
	if err != nil {
		return nil, err
	}

	return o.Id, nil
}
