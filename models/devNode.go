package models

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/Unknwon/com"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	"github.com/json-iterator/go"
	_ "github.com/mattn/go-sqlite3"
	"github.com/slzm40/common"
)

const (
	_DB_NAME              = "data/devll.db"
	_DB_DRIVE             = "sqlite3"
	_default_trunkid_list = `{"trunkID":[]}`
	_default_bind_list    = `{"id":[]}`
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
	ID           uint   //`gorm:"primary_key"`
	IeeeAddr     uint64 //`gorm:"NOT NULL"`
	NodeNo       uint16 //`gorm:"UNIQUE;NOT NULL"`
	NwkAddr      uint16 //`gorm:"UNIQUE;NOT NULL"`
	InTrunkList  string // 输入集表
	OutTrunkList string // 输出集表
	SrcBindList  string // 源绑定表 : 谁绑定了本设备
	DstBindList  string // 目的绑定表: 本设备绑定了谁
	Status       uint32 // 状态掩码 保留
	inTrunk      []uint16
	outTrunk     []uint16
	srcBind      []uint
	dstBind      []uint
}

const deviceNodeInfos_Sql = `CREATE TABLE "device_node_infos" (
	"id" integer primary key autoincrement,
	"ieee_addr" bigint NOT NULL,
	"node_no" integer NOT NULL,
	"nwk_addr" integer NOT NULL,
	"in_trunk_list" varchar(255),
	"out_trunk_list" varchar(255),
	"src_bind_list" varchar(255),
	"dst_bind_list" varchar(255),
	"status" integer default(0),
	UNIQUE(ieee_addr,node_no) ON CONFLICT FAIL)`

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
	//devDb.LogMode(true)

	devDb.AutoMigrate(&DeviceInfo{})
	if devDb.Error != nil {
		panic("models: gorm AutoMigrate failed," + err.Error())
	}

	if !devDb.HasTable("device_node_infos") {
		devDb.Raw(deviceNodeInfos_Sql).Scan(&DeviceNodeInfo{})
	}

}

// 根据网络地址,节点号找到设备节点
func LookupDeviceNodeByNN(nwkAddr, nodeNum uint16) (*DeviceNodeInfo, error) {
	o := &DeviceNodeInfo{}
	if devDb.Where(&DeviceNodeInfo{NwkAddr: nwkAddr, NodeNo: nodeNum}).First(o).RecordNotFound() == true {
		return nil, gorm.ErrRecordNotFound
	}

	if err := o.parseInternalJsonString(); err != nil {
		return nil, err
	}

	return o, nil
}

// 根据ieee地址,节点号找到设备节点
func LookupDeviceNodeByIN(ieeeAddr uint64, nodeNum uint16) (*DeviceNodeInfo, error) {
	o := &DeviceNodeInfo{}
	if devDb.Where(&DeviceNodeInfo{IeeeAddr: ieeeAddr, NodeNo: nodeNum}).First(o).RecordNotFound() == true {
		return nil, gorm.ErrRecordNotFound
	}

	if err := o.parseInternalJsonString(); err != nil {
		return nil, err
	}

	return o, nil
}

// 根据id找到设备节点
func LookupDeviceNodeByID(id uint) (*DeviceNodeInfo, error) {
	o := &DeviceNodeInfo{}
	if devDb.First(o, id).RecordNotFound() == true {
		return nil, gorm.ErrRecordNotFound
	}

	if err := o.parseInternalJsonString(); err != nil {
		return nil, err
	}

	return o, nil
}

// 根据id找到设备节点
func lookupDeviceNodeByID(db *gorm.DB, id uint) (*DeviceNodeInfo, error) {
	o := &DeviceNodeInfo{}
	if db.First(o, id).RecordNotFound() == true {
		return nil, gorm.ErrRecordNotFound
	}

	if err := o.parseInternalJsonString(); err != nil {
		return nil, err
	}

	return o, nil
}

// 绑定两个设备 要求更新 源设备节点的<目的绑定表>和目标设备节点的<源绑定表>
func BindDeviceNode(SrcIeeeAddr uint64, SrcNodeNum uint16, DstBindIeeeAddr uint64, DstBindNodeNum, BindTrunkID uint16) error {
	var (
		SrcDNI, DstDNI *DeviceNodeInfo
		err            error
	)
	// 获取源设备节点和目的设备节点
	if SrcDNI, err = LookupDeviceNodeByIN(SrcIeeeAddr, SrcNodeNum); err != nil {
		return err
	}
	if DstDNI, err = LookupDeviceNodeByIN(DstBindIeeeAddr, DstBindNodeNum); err != nil {
		return err
	}

	// 只有源设备节点输出集和目的设备输入集互补,即都含有要绑定的集,才进行绑定
	if !common.IsSliceContainsUint16(SrcDNI.outTrunk, BindTrunkID) || !common.IsSliceContainsUint16(DstDNI.inTrunk, BindTrunkID) {
		return errors.New("src and dst trunkID 不是互补")
	}

	// 源设备节点 目的绑定表 不含目标设备节点 或 目标设备节点 源绑定表 不含源设备节点 将进行绑定添加, 都有直接返回成功
	if common.IsSliceContainsUint(SrcDNI.dstBind, DstDNI.ID) && common.IsSliceContainsUint(DstDNI.srcBind, SrcDNI.ID) {
		return nil
	}

	SrcDNI_DstBd := SrcDNI.dstBind
	DstDNI_SrcBd := DstDNI.srcBind

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

	// 更新本地绑定
	SrcDNI.dstBind = SrcDNI_DstBd
	DstDNI.srcBind = DstDNI_SrcBd

	return nil
}

// 解绑两个设备节点
// 如果两个设备节点有两个互补集绑定的,那么将进行同时解绑
func UnBindDeviceNode(SrcIeeeAddr uint64, SrcNodeNum uint16, DstBindIeeeAddr uint64, DstBindNodeNum, BindTrunkID uint16) error {
	var (
		SrcDNI, DstDNI *DeviceNodeInfo
		err            error
	)
	// 获取源设备节点和目的设备节点
	if SrcDNI, err = LookupDeviceNodeByIN(SrcIeeeAddr, SrcNodeNum); err != nil {
		return nil
	}
	if DstDNI, err = LookupDeviceNodeByIN(DstBindIeeeAddr, DstBindNodeNum); err != nil {
		return nil
	}

	// 只有源设备节点输出集和目的设备输入集互补,即都含有要绑定的集,才进行解绑定,否则认为是成功的
	if !common.IsSliceContainsUint16(SrcDNI.outTrunk, BindTrunkID) || !common.IsSliceContainsUint16(DstDNI.inTrunk, BindTrunkID) {
		return nil
	}

	// 源设备节点的<目的绑定表>不含目标设备节点
	//或 目标设备节点<源绑定表>不含源设备节点 将进行绑定解绑直接返回成功
	if !common.IsSliceContainsUint(SrcDNI.dstBind, DstDNI.ID) || !common.IsSliceContainsUint(DstDNI.srcBind, SrcDNI.ID) {
		return nil
	}

	SrcDNI_DstBd := SrcDNI.dstBind
	DstDNI_SrcBd := DstDNI.srcBind

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
	// 更新本地
	SrcDNI.dstBind = SrcDNI_DstBd
	DstDNI.srcBind = DstDNI_SrcBd

	return nil
}

// 找到绑定表的所有设备节点
func BindFindDeviceNodeByNN(NwkAddr, NodeNum, trunkID uint16) ([]*DeviceNodeInfo, error) {
	src, err := LookupDeviceNodeByNN(NwkAddr, NodeNum)
	if err != nil {
		return nil, err
	}

	// 源设备 是否包含输出集
	if common.IsSliceContainsUint16(src.outTrunk, trunkID) != true {
		return nil, errors.New("不包含对应的集")
	}

	dni := make([]*DeviceNodeInfo, 0, len(src.dstBind))
	for _, id := range src.dstBind {
		tmpdni, err := LookupDeviceNodeByID(uint(id))
		if err != nil {
			continue
		}

		// 只有目标设备含有输入集才加入
		if common.IsSliceContainsUint16(tmpdni.inTrunk, trunkID) {
			dni = append(dni, tmpdni)
		}
	}

	return dni, nil
}

// 找到绑定表的所有设备节点
func BindFindDeviceNodeByIN(ieeeAddr uint64, NodeNum, trunkID uint16) ([]*DeviceNodeInfo, error) {
	src, err := LookupDeviceNodeByIN(ieeeAddr, NodeNum)
	if err != nil {
		return nil, err
	}

	// 源设备 是否包含输出集
	if common.IsSliceContainsUint16(src.outTrunk, trunkID) != true {
		return nil, errors.New("不包含对应的集")
	}

	dni := make([]*DeviceNodeInfo, 0, len(src.dstBind))
	for _, id := range src.dstBind {
		tmpdni, err := LookupDeviceNodeByID(uint(id))
		if err != nil {
			continue
		}

		// 只有目标设备含有输入集才加入
		if common.IsSliceContainsUint16(tmpdni.inTrunk, trunkID) {
			dni = append(dni, tmpdni)
		}
	}

	return dni, nil
}

// 解析内部的jsonString
func (this *DeviceNodeInfo) parseInternalJsonString() error {
	var err error

	// 输入集
	tmpInTk := &TrunkIDList{}
	if err = jsoniter.UnmarshalFromString(this.InTrunkList, tmpInTk); err != nil {
		return err
	}
	// 输出集
	tmpOutTk := &TrunkIDList{}
	if err = jsoniter.UnmarshalFromString(this.OutTrunkList, tmpOutTk); err != nil {
		return err
	}

	// 源绑定
	oSrc := &BindInfo{}
	if err = jsoniter.UnmarshalFromString(this.SrcBindList, oSrc); err != nil {
		return err
	}
	// 目标绑定
	oDst := &BindInfo{}
	if err = jsoniter.UnmarshalFromString(this.DstBindList, oDst); err != nil {
		return err
	}

	this.inTrunk = tmpInTk.TrunkID
	this.outTrunk = tmpOutTk.TrunkID
	this.srcBind = oSrc.Id
	this.dstBind = oDst.Id

	return nil
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
func (this *DeviceNodeInfo) GetTrunkIDList() (inTrunk, outTrunk []uint16) {
	return this.inTrunk, this.outTrunk
}

// 设置设备节点 集id表
func (this *DeviceNodeInfo) SetTrunkIDlist(inTrunk, outTrunk []uint16) error {
	var err error

	if len(inTrunk) == 0 {
		this.InTrunkList = _default_trunkid_list
	} else if this.InTrunkList, err = jsoniter.MarshalToString(&TrunkIDList{TrunkID: inTrunk}); err != nil {
		return err
	}

	if len(outTrunk) == 0 {
		this.OutTrunkList = _default_trunkid_list
	} else if this.OutTrunkList, err = jsoniter.MarshalToString(&TrunkIDList{TrunkID: outTrunk}); err != nil {
		return err
	}

	return nil
}

// 获取设备节点源绑定id列表
func (this *DeviceNodeInfo) GetSrcBindList() []uint {
	return this.srcBind
}

// 设置设备节点源绑定id列表
func (this *DeviceNodeInfo) setSrcBindList(id []uint) error {
	var err error

	if len(id) == 0 {
		this.SrcBindList = _default_bind_list
	} else if this.SrcBindList, err = jsoniter.MarshalToString(&BindInfo{Id: id}); err != nil {
		return err
	}

	return nil
}

// 获取设备节点目的绑定id列表
func (this *DeviceNodeInfo) GetDstBindList() []uint {
	return this.dstBind
}

// 设置设备节点目的绑定id列表
func (this *DeviceNodeInfo) setDstBindList(id []uint) error {
	var err error

	if len(id) == 0 {
		this.DstBindList = _default_bind_list
	} else if this.DstBindList, err = jsoniter.MarshalToString(&BindInfo{Id: id}); err != nil {
		return err
	}

	return nil
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
	if newCapacityValue == this.Capacity {
		return nil
	}

	devDb.Model(this).Update("capacity", newCapacityValue)
	return devDb.Error
}

// 更新设备和设备节点所有节点的网络地址
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

	//创建除保留节点(0)外的所有节点
	for i, v := range devNode {
		dnode := &DeviceNodeInfo{
			IeeeAddr: this.IeeeAddr,
			NodeNo:   uint16(i + 1),
			NwkAddr:  this.NwkAddr,
		}
		if err = dnode.SetTrunkIDlist(v.InTrunk, v.OutTrunk); err != nil {
			dnode.InTrunkList = _default_trunkid_list
			dnode.OutTrunkList = _default_trunkid_list
			logs.Warning("models: SetTrunkIDlist ", err)
		}

		dnode.SrcBindList = _default_bind_list
		dnode.DstBindList = _default_bind_list

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
		fmt.Println("bug1")
		return tx.Error
	}

	var devNodes []DeviceNodeInfo
	tx.Find(&devNodes, &DeviceNodeInfo{IeeeAddr: this.IeeeAddr})

	for _, v := range devNodes {
		// 是否有源绑定,有则删除每一个 源的目标绑定
		if len(v.srcBind) > 0 {
			for _, tv := range v.srcBind { // 扫描每一个源 的目标绑定,让它删除对应id
				tmpdevNode, err := lookupDeviceNodeByID(tx, uint(tv))
				if err == nil && len(tmpdevNode.dstBind) > 0 && common.IsSliceContainsUint(tmpdevNode.dstBind, v.ID) {
					tmpdevNode.dstBind = common.DeleteFromSliceUint(tmpdevNode.dstBind, v.ID)

					if err = tmpdevNode.setDstBindList(tmpdevNode.dstBind); err != nil {
						continue
					}

					tx.Model(tmpdevNode).Updates(&DeviceNodeInfo{DstBindList: tmpdevNode.DstBindList})
					if tx.Error != nil {
						tx.Rollback()
						fmt.Println("bug2")
						return tx.Error
					}
				}
			}
		}

		// 是否有目标绑定, 有则删除每一个 目标的源绑定
		if len(v.dstBind) > 0 {
			for _, tv := range v.srcBind { // 扫描每一个源 的目标绑定,让它删除对应id
				tmpdevNode, err := lookupDeviceNodeByID(tx, uint(tv))
				if err == nil && len(tmpdevNode.srcBind) > 0 && common.IsSliceContainsUint(tmpdevNode.srcBind, v.ID) {
					tmpdevNode.srcBind = common.DeleteFromSliceUint(tmpdevNode.srcBind, v.ID)

					if err = tmpdevNode.setDstBindList(tmpdevNode.srcBind); err != nil {
						continue
					}

					tx.Model(tmpdevNode).Updates(&DeviceNodeInfo{SrcBindList: tmpdevNode.SrcBindList})
					if tx.Error != nil {
						tx.Rollback()
						fmt.Println("bug3")
						return tx.Error
					}
				}
			}
		}

		tx.Unscoped().Delete(v)
		if tx.Error != nil {
			tx.Rollback()
			fmt.Println("bug4")
			return tx.Error
		}
	}

	if err := tx.Unscoped().Delete(this).Error; err != nil {
		tx.Rollback()
		fmt.Println("bug5")
		return err
	}

	if tx.Commit().Error != nil {
		tx.Rollback()
		return tx.Error
	}

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
		//删除设备和设备所有的节点
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

	// 更新能力属性
	if dev.Capacity != capacity {
		if err = dev.updateCapacity(capacity); err != nil {
			return err
		}
	}

	// 更新对应网络址
	if dev.NwkAddr != nwkAddr {
		if err := dev.updateDeviceAndNodeNwkAddr(nwkAddr); err != nil {
			return err
		}
	}

	return nil
}

//删除这个设备及设备的所有节点
func DeleteDeveiceAndNode(ieeeAddr uint64) error {
	dev, err := LookupDeviceByIeeeAddr(ieeeAddr)
	if err != nil {
		return nil
	}

	return dev.deleteDeveiceAndNode()
}
