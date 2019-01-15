package models

import (
	"errors"
	"os"
	"path"

	"github.com/astaxie/beego/logs"

	"github.com/slzm40/common"

	"github.com/Unknwon/com"
	"github.com/jinzhu/gorm"
	"github.com/json-iterator/go"
	_ "github.com/mattn/go-sqlite3"
)

const (
	_DB_NAME              = "data/devll.db"
	_DB_DRIVE             = "sqlite3"
	_DEFAULT_TRUNKID_LIST = `{"trunkID":[]}`
	_DEFAULT_BIND_LIST    = `{"id":[]}`
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
type DeviceNodeInfo struct {
	ID           uint   `gorm:"primary_key"`
	NwkAddr      uint16 `gorm:"UNIQUE;NOT NULL"`
	NodeNo       uint16 `gorm:"UNIQUE;NOT NULL"`
	IeeeAddr     uint64 `gorm:"NOT NULL"`
	InTrunkList  string `gorm:'default:"{\"trunkID\":[]}"'`
	OutTrunkList string `gorm:'default:"{\"trunkID\":[]}"'`
	SrcBindList  string `gorm:'default:"{\"id\":[]}'` // 源绑定表 : 谁绑定了本设备
	DstBindList  string `gorm:'default:"{\"id\":[]}'` // 目的绑定表: 本设备绑定了谁
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

	devDb.AutoMigrate(&Product{}, &DeviceInfo{}, &DeviceNodeInfo{})
}

// 根据网络地址,节点号找到设备节点
func LookupDeviceNodeByNN(nwkAddr, nodeNum uint16) (*DeviceNodeInfo, error) {
	o := &DeviceNodeInfo{}
	if devDb.Where(&DeviceNodeInfo{NwkAddr: nwkAddr, NodeNo: nodeNum}).First(o).RecordNotFound() == true {
		return nil, gorm.ErrRecordNotFound
	}

	return o, nil
}

// 根据ieee地址,节点号找到设备节点
func LookupDeviceNodeByIN(ieeeAddr uint64, nodeNum uint16) (*DeviceNodeInfo, error) {
	o := &DeviceNodeInfo{}
	if devDb.Where(&DeviceNodeInfo{IeeeAddr: ieeeAddr, NodeNo: nodeNum}).First(o).RecordNotFound() == true {
		return nil, gorm.ErrRecordNotFound
	}

	return o, nil
}

// 根据id找到设备节点
func LookupDeviceNodeByID(id uint) (*DeviceNodeInfo, error) {
	o := &DeviceNodeInfo{}
	if devDb.First(o, id).RecordNotFound() == true {
		return nil, gorm.ErrRecordNotFound
	}

	return o, nil
}

// 绑定两个设备 要求更新 源设备节点的 目的绑定表和目标设备节点的 源绑定表
func BindDeviceNode(SrcIeeeAddr, DstBindIeeeAddr uint64, SrcNodeNum, DstBindNodeNum, BindTrunkID uint16) error {
	var (
		SrcDNI, DstDNI             *DeviceNodeInfo
		SrcOutTk, DstInTk          []uint16
		SrcDNI_DstBd, DstDNI_SrcBd []uint
		err                        error
	)
	// 获取源设备节点和目的设备节点
	if SrcDNI, err = LookupDeviceNodeByIN(SrcIeeeAddr, SrcNodeNum); err != nil {
		return err
	}
	if DstDNI, err = LookupDeviceNodeByIN(DstBindIeeeAddr, DstBindNodeNum); err != nil {
		return err
	}
	//获取源设备节点输出集和目的设备输入集
	if _, SrcOutTk, err = SrcDNI.GetTrunkIDList(); err != nil {
		return err
	}
	if DstInTk, _, err = DstDNI.GetTrunkIDList(); err != nil {
		return err
	}
	// 只有源设备节点输出集和目的设备输入集互补,即都含有要绑定的设备,才进行绑定
	if !common.IsSliceContainsUint16(SrcOutTk, BindTrunkID) || !common.IsSliceContainsUint16(DstInTk, BindTrunkID) {
		return errors.New("src and dst trunkID 不是互补")
	}

	//获取源设备节点-目的绑定表 和目标设备节点-源绑定表
	if SrcDNI_DstBd, err = SrcDNI.GetDstBindList(); err != nil {
		return err
	}
	if DstDNI_SrcBd, err = DstDNI.GetSrcBindList(); err != nil {
		return err
	}

	// 源设备节点 目的绑定表 不含目标设备节点 或 目标设备节点 源绑定表 不含源设备节点 将进行绑定添加, 都有直接返回成功
	if common.IsSliceContainsUint(SrcDNI_DstBd, DstDNI.ID) && common.IsSliceContainsUint(DstDNI_SrcBd, SrcDNI.ID) {
		return nil
	}
	SrcDNI_DstBd = common.AppendUint(SrcDNI_DstBd, DstDNI.ID)
	DstDNI_SrcBd = common.AppendUint(DstDNI_SrcBd, SrcDNI.ID)

	if err = SrcDNI.setDstBindList(SrcDNI_DstBd); err != nil {
		return err
	}
	if err = DstDNI.setSrcBindList(DstDNI_SrcBd); err != nil {
		return err
	}
	// 开始更新表
	tx := devDb.Begin()
	if err = tx.Error; err != nil {
		return err
	}

	tx.Model(&SrcDNI).Updates(&DeviceNodeInfo{DstBindList: SrcDNI.DstBindList})
	if err = tx.Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Model(&DstDNI).Updates(&DeviceNodeInfo{SrcBindList: DstDNI.SrcBindList})
	if err = tx.Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	if err = tx.Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// 解绑两个设备节点
// 如果两个设备节点有两个互补集进行绑定,那么将进行同时解绑
func UnBindDeviceNode(SrcIeeeAddr, DstBindIeeeAddr uint64, SrcNodeNum, DstBindNodeNum, BindTrunkID uint16) error {
	var (
		SrcDNI, DstDNI             *DeviceNodeInfo
		SrcOutTk, DstInTk          []uint16
		SrcDNI_DstBd, DstDNI_SrcBd []uint
		err                        error
	)
	// 获取源设备节点和目的设备节点
	if SrcDNI, err = LookupDeviceNodeByIN(SrcIeeeAddr, SrcNodeNum); err != nil {
		return nil
	}
	if DstDNI, err = LookupDeviceNodeByIN(DstBindIeeeAddr, DstBindNodeNum); err != nil {
		return nil
	}
	//获取源设备节点输出集和目的设备输入集
	if _, SrcOutTk, err = SrcDNI.GetTrunkIDList(); err != nil {
		return err
	}
	if DstInTk, _, err = DstDNI.GetTrunkIDList(); err != nil {
		return err
	}
	// 只有源设备节点输出集和目的设备输入集互补,即都含有要绑定的集,才进行解绑定,否则认为是成功的
	if !common.IsSliceContainsUint16(SrcOutTk, BindTrunkID) || !common.IsSliceContainsUint16(DstInTk, BindTrunkID) {
		return nil
	}

	//获取源设备节点-目的绑定表 和目标设备节点-源绑定表
	if SrcDNI_DstBd, err = SrcDNI.GetDstBindList(); err != nil {
		return err
	}
	if DstDNI_SrcBd, err = DstDNI.GetSrcBindList(); err != nil {
		return err
	}

	// 源设备节点 目的绑定表 不含目标设备节点 或 目标设备节点 源绑定表 不含源设备节点 将进行绑定添加, 都有直接返回成功
	if common.IsSliceContainsUint(SrcDNI_DstBd, DstDNI.ID) && common.IsSliceContainsUint(DstDNI_SrcBd, SrcDNI.ID) {
		return nil
	}
	SrcDNI_DstBd = common.DeleteFromSliceUintAll(SrcDNI_DstBd, DstDNI.ID)
	DstDNI_SrcBd = common.DeleteFromSliceUintAll(DstDNI_SrcBd, SrcDNI.ID)

	if err = SrcDNI.setDstBindList(SrcDNI_DstBd); err != nil {
		return err
	}
	if err = DstDNI.setSrcBindList(DstDNI_SrcBd); err != nil {
		return err
	}
	// 开始更新表
	tx := devDb.Begin()
	if err = tx.Error; err != nil {
		return err
	}

	tx.Model(&SrcDNI).Updates(&DeviceNodeInfo{DstBindList: SrcDNI.DstBindList})
	if err = tx.Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Model(&DstDNI).Updates(&DeviceNodeInfo{SrcBindList: DstDNI.SrcBindList})
	if err = tx.Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	if err = tx.Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// 找到绑定表的所有设备节点
func BindFindDeviceNode(srcNwkAddr, srcNodeNum, trunkID uint16) ([]*DeviceNodeInfo, error) {
	src, err := LookupDeviceNodeByNN(srcNwkAddr, srcNodeNum)
	if err != nil {
		return nil, err
	}

	_, outTK, err := src.GetTrunkIDList()
	if err != nil {
		return nil, err
	}
	if common.IsSliceContainsUint16(outTK, trunkID) != true {
		return nil, err
	}

	bindId, err := src.GetDstBindList()
	if err != nil {
		return nil, err
	}

	dni := make([]*DeviceNodeInfo, 0, len(bindId))
	for id := range bindId {
		tmpdni, err := LookupDeviceNodeByID(uint(id))
		if err != nil {
			continue
		}

		inTK, _, err := tmpdni.GetTrunkIDList()
		if err != nil {
			continue
		}

		if common.IsSliceContainsUint16(inTK, trunkID) {
			dni = append(dni, tmpdni)
		}
	}

	return dni, nil
}

// 获取设备节点id
func (this *DeviceNodeInfo) GetID() uint {
	return this.ID
}

// 获取设备节点网络地址
func (this *DeviceNodeInfo) GetNwkAddr() uint16 {
	return this.NwkAddr
}

// 获取设备节点节点号
func (this *DeviceNodeInfo) GetNodeNum() uint16 {
	return this.NodeNo
}

// 获取设备节点Ieee地址
func (this *DeviceNodeInfo) GetIeeeAddr() uint64 {
	return this.IeeeAddr
}

// 获取设备节点 集id表
func (this *DeviceNodeInfo) GetTrunkIDList() (inTrunk, outTrunk []uint16, err error) {
	if len(this.InTrunkList) == 0 || len(this.OutTrunkList) == 0 {
		return nil, nil, errors.New("device node info trunklist nil")
	}

	tmpInTk := &TrunkIDList{}
	if err = jsoniter.UnmarshalFromString(this.InTrunkList, tmpInTk); err != nil {
		return nil, nil, err
	}

	tmpOutTk := &TrunkIDList{}
	if err = jsoniter.UnmarshalFromString(this.OutTrunkList, tmpOutTk); err != nil {
		return nil, nil, err
	}

	return tmpInTk.TrunkID, tmpOutTk.TrunkID, nil
}

// 设置设备节点 集id表
func (this *DeviceNodeInfo) SetTrunkIDlist(inTrunk, outTrunk []uint16) error {
	var err error

	tmpTk := &TrunkIDList{}

	if len(inTrunk) != 0 {
		tmpTk.TrunkID = inTrunk
	} else {
		tmpTk.TrunkID = []uint16{}
	}
	if this.InTrunkList, err = jsoniter.MarshalToString(tmpTk); err != nil {
		return err
	}

	if len(outTrunk) != 0 {
		tmpTk.TrunkID = outTrunk
	} else {
		tmpTk.TrunkID = []uint16{}
	}
	if this.OutTrunkList, err = jsoniter.MarshalToString(tmpTk); err != nil {
		return err
	}

	return nil
}

// 获取设备节点目的绑定id列表
func (this *DeviceNodeInfo) GetDstBindList() ([]uint, error) {
	o := &BindInfo{}
	err := jsoniter.UnmarshalFromString(this.DstBindList, o)
	if err != nil {
		return nil, err
	}

	return o.Id, nil
}

// 设置设备节点目的绑定id列表
func (this *DeviceNodeInfo) setDstBindList(id []uint) error {
	var err error

	this.DstBindList, err = jsoniter.MarshalToString(&BindInfo{Id: id})
	if err != nil {
		return err
	}

	return nil
}

// 获取设备节点源绑定id列表
func (this *DeviceNodeInfo) GetSrcBindList() ([]uint, error) {
	o := &BindInfo{}
	err := jsoniter.UnmarshalFromString(this.SrcBindList, o)
	if err != nil {
		return nil, err
	}

	return o.Id, nil
}

// 设置设备节点源绑定id列表
func (this *DeviceNodeInfo) setSrcBindList(id []uint) error {
	var err error

	this.SrcBindList, err = jsoniter.MarshalToString(&BindInfo{Id: id})
	if err != nil {
		return err
	}

	return nil
}

// 设备节点增加一个绑定id,如果是新增返回true,原来列表就有返回false
func (this *DeviceNodeInfo) AddDstBindList(id uint) (isNewAdd bool, err error) {
	ids, err := this.GetDstBindList()
	if err != nil {
		return false, err
	}

	if common.IsSliceContainsUint(ids, id) {
		return false, nil
	}
	ids = append(ids, id)

	return true, this.setDstBindList(ids)
}

// 设备节点删除一个绑定id,如果是新删返回true,原来列表就没有返回false
func (this *DeviceNodeInfo) DeleteDstBindList(id uint) (isNewDel bool, err error) {
	ids, err := this.GetDstBindList()
	if err != nil {
		return false, err
	}

	for i, v := range ids {
		if v == id {
			s := append(ids[0:i], ids[i+1:]...)
			return true, this.setDstBindList(s)
		}
	}

	return false, nil
}

// 根据网络地址找到设备
func LookupDeviceByNwkAddr(nwkAddr uint16) (*DeviceInfo, error) {
	oInfo := &DeviceInfo{}
	if db := devDb.Where(&DeviceInfo{NwkAddr: nwkAddr}).First(oInfo); db.RecordNotFound() == true {
		return nil, gorm.ErrRecordNotFound
	}

	return oInfo, nil
}

// 根据ieee地址找到设备
func LookupDeviceByIeeeAddr(ieeeAddr uint64) (*DeviceInfo, error) {
	oInfo := &DeviceInfo{}
	if db := devDb.Where(&DeviceInfo{IeeeAddr: ieeeAddr}).First(oInfo); db.RecordNotFound() == true {
		return nil, gorm.ErrRecordNotFound
	}

	return oInfo, nil
}

// 获取ieee地址
func (this *DeviceInfo) GetIeeeAddr() uint64 {
	return this.IeeeAddr
}

// 获取网络地址
func (this *DeviceInfo) GetNwkAddr() uint16 {
	return this.NwkAddr
}

// 获取设备能力
func (this *DeviceInfo) GetCapacity() byte {
	return this.Capacity
}

// 获取设备的id
func (this *DeviceInfo) GetID() uint {
	return this.ID
}

// 获取设备的产品id
func (this *DeviceInfo) GetProductID() uint32 {
	return this.ProductId
}

func (this *DeviceInfo) updateCapacity(newCapacityValue byte) error {
	return nil
}

// 更新设备和设备所有节点的网络地址
func (this *DeviceInfo) updateDeviceAndNodeNwkAddr(NewnwkAddr uint16) error {
	var err error

	tx := devDb.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	// 更新设备网络地址
	if err = tx.Model(this).Updates(&DeviceInfo{NwkAddr: NewnwkAddr}).Error; err != nil {
		tx.Rollback()
		return err
	}
	//更新所有节点网络地址
	if err = tx.Model(&DeviceNodeInfo{}).Where(&DeviceNodeInfo{IeeeAddr: this.IeeeAddr}).Updates(&DeviceNodeInfo{NwkAddr: NewnwkAddr}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// 创建设备和设备所有的节点,失败将不建立
func (this *DeviceInfo) createDeveiceAndNode() error {
	// 查询对应产品
	pdt, err := LookupProduct(this.ProductId)
	if err != nil {
		return err
	}

	devNode, err := pdt.GetDeviceNodeDscList()
	if err != nil {
		return err
	}

	tx := devDb.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	// 创建设备
	tx.Create(this)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	//创建除保留节点0外的所有节点
	for i, v := range devNode {
		dnode := &DeviceNodeInfo{
			NwkAddr:  this.NwkAddr,
			NodeNo:   uint16(i + 1),
			IeeeAddr: this.IeeeAddr,
		}
		if err = dnode.SetTrunkIDlist(v.InTrunk, v.OutTrunk); err != nil {
			dnode.InTrunkList = _DEFAULT_TRUNKID_LIST
			dnode.OutTrunkList = _DEFAULT_TRUNKID_LIST
			logs.Warning("models: SetTrunkIDlist ", err)
		}

		if err = tx.Create(dnode).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// 根据ieee地址删除设备,包含所有的设备节点
func (this *DeviceInfo) deleteDeveiceAndNode() error {
	tx := devDb.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Unscoped().Delete(this).Error; err != nil {
		tx.Rollback()
		return err
	}

	// TODO:删除所有节点
	return nil
}

// 添加设备信息,包括所有设备节点,如果已存在,则更新相应参数
func UpdateDeviceAndANode(ieeeAddr uint64, nwkAddr uint16, capacity byte, productID uint32) error {
	var err error

	dev, err := LookupDeviceByIeeeAddr(ieeeAddr)
	if err != nil {
		return (&DeviceInfo{
			IeeeAddr:  ieeeAddr,
			NwkAddr:   nwkAddr,
			Capacity:  capacity,
			ProductId: productID,
		}).createDeveiceAndNode()
	}

	// 已存在,看是否要更新相应参数
	if dev.ProductId != productID {
		//删除所有设备和设备节点
		if err = dev.deleteDeveiceAndNode(); err != nil {
			return err
		}

		//重建所有设备和节点
		return (&DeviceInfo{
			IeeeAddr:  ieeeAddr,
			NwkAddr:   nwkAddr,
			Capacity:  capacity,
			ProductId: productID,
		}).createDeveiceAndNode()
	}

	if dev.Capacity != capacity {
		if err = dev.updateCapacity(capacity); err != nil {
			return err
		}
	}

	if dev.NwkAddr != nwkAddr {
		if err := dev.updateDeviceAndNodeNwkAddr(nwkAddr); err != nil {
			return err
		}
	}

	return nil
}

func DeleteDeveiceAndNode(ieeeAddr uint64) error {
	dev, err := LookupDeviceByIeeeAddr(ieeeAddr)
	if err != nil {
		return nil
	}

	return dev.deleteDeveiceAndNode()
}
