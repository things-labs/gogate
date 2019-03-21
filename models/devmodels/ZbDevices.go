package devmodels

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/slzm40/common"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

// 设备表
type ZbDeviceInfo struct {
	gorm.Model
	IeeeAddr  uint64 `gorm:"UNIQUE;NOT NULL"`
	NwkAddr   uint16 `gorm:"NOT NULL"`
	Capacity  byte   `gorm:"default:2"`
	ProductId int
}

//NOTE: 表内数组以逗号","分隔
// 节点表
type ZbDeviceNodeInfo struct {
	ID           uint   //`gorm:"primary_key"`
	IeeeAddr     uint64 //`gorm:"NOT NULL"`
	NodeNo       uint16 //`gorm:"UNIQUE;NOT NULL"`
	NwkAddr      uint16 //`gorm:"UNIQUE;NOT NULL"`
	ProductId    int
	InTrunkList  string // 输入集表
	OutTrunkList string // 输出集表
	SrcBindList  string // 源绑定表 : 谁绑定了本设备
	DstBindList  string // 目的绑定表: 本设备绑定了谁
	Status       uint32 // 状态掩码 保留
	// 以下4个用于即时解析处理,不在数据库
	inTrunk  []string // 输入集表
	outTrunk []string // 输出集表
	srcBind  []string // 源绑定表 : 谁绑定了本设备
	dstBind  []string // 目的绑定表: 本设备绑定了谁
}

const zbDeviceNodeInfos_Sql = `CREATE TABLE "zb_device_node_infos" (
	"id" integer primary key autoincrement,
	"ieee_addr" bigint NOT NULL,
	"node_no" integer NOT NULL,
	"nwk_addr" integer NOT NULL,
	"product_id" integer,
	"in_trunk_list" varchar(255),
	"out_trunk_list" varchar(255),
	"src_bind_list" varchar(255),
	"dst_bind_list" varchar(255),
	"status" integer default(0),
	UNIQUE(ieee_addr,node_no) ON CONFLICT FAIL)`

// 采用","分隔字符串
func splitInternalString(s string) []string {
	if s == "" {
		return []string{}
	}

	return strings.Split(s, ",") // NOTE: 如果切割空字符串会返回一个空字符的数组,长度为1,这是个坑
}

// 采用","拼接字符串
func joinInternalString(s []string) string {
	return strings.Join(s, ",")
}

// 根据网络地址,节点号找到设备节点
func LookupZbDeviceNodeByNN(nwkAddr, nodeNum uint16) (*ZbDeviceNodeInfo, error) {
	o := &ZbDeviceNodeInfo{}
	if devDb.Where(&ZbDeviceNodeInfo{
		NwkAddr: nwkAddr,
		NodeNo:  nodeNum}).First(o).RecordNotFound() {
		return nil, gorm.ErrRecordNotFound
	}

	o.parseInternalString()

	return o, nil
}

// 根据ieee地址,节点号找到设备节点
func LookupZbDeviceNodeByIN(ieeeAddr uint64, nodeNum uint16) (*ZbDeviceNodeInfo, error) {
	o := &ZbDeviceNodeInfo{}
	if devDb.Where(&ZbDeviceNodeInfo{
		IeeeAddr: ieeeAddr,
		NodeNo:   nodeNum}).First(o).RecordNotFound() {
		return nil, gorm.ErrRecordNotFound
	}

	o.parseInternalString()

	return o, nil
}

// 根据id找到设备节点
func LookupZbDeviceNodeByID(id string) (*ZbDeviceNodeInfo, error) {
	idval, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		return nil, err
	}

	o := &ZbDeviceNodeInfo{}
	if devDb.First(o, int(idval)).RecordNotFound() {
		return nil, gorm.ErrRecordNotFound
	}

	o.parseInternalString()

	return o, nil
}

// 根据id找到设备节点
func lookupZbDeviceNodeByID(db *gorm.DB, id string) (*ZbDeviceNodeInfo, error) {
	idval, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		return nil, err
	}
	o := &ZbDeviceNodeInfo{}
	if db.First(o, idval).RecordNotFound() {
		return nil, gorm.ErrRecordNotFound
	}

	o.parseInternalString()

	return o, nil
}

func HasZbDevice(ieeeAddr uint64, pid int) bool {
	if !HasZbProduct(pid) {
		return false
	}
	if db := devDb.Where(&ZbDeviceInfo{IeeeAddr: ieeeAddr}).First(&ZbDeviceInfo{}); db.RecordNotFound() {
		return false
	}
	return true
}

// 绑定两个设备 要求更新 源设备节点的<目的绑定表>和目标设备节点的<源绑定表>
func BindZbDeviceNode(SrcIeeeAddr uint64, SrcNodeNum uint16,
	DstBindIeeeAddr uint64, DstBindNodeNum, BindTrunkID uint16) error {
	var SrcDNI *ZbDeviceNodeInfo
	var DstDNI *ZbDeviceNodeInfo
	var err error

	// 获取源设备节点
	if SrcDNI, err = LookupZbDeviceNodeByIN(SrcIeeeAddr, SrcNodeNum); err != nil {
		return err
	}
	// 获取目的设备节点
	if DstDNI, err = LookupZbDeviceNodeByIN(DstBindIeeeAddr, DstBindNodeNum); err != nil {
		return err
	}

	strBindTkid := common.FormatBaseTypes(BindTrunkID)
	// 只有源设备节点输出集和目的设备输入集互补,即都含有要绑定的集,才进行绑定
	if !common.IsSliceContainsStrNocase(SrcDNI.outTrunk, strBindTkid) ||
		!common.IsSliceContainsStrNocase(DstDNI.inTrunk, strBindTkid) {
		return errors.New("src and dst trunkID Not a complementary")
	}

	// 源设备节点 目的绑定表 不含目标设备节点 或 目标设备节点 源绑定表 不含源设备节点 将进行绑定添加,
	// 都有直接返回成功
	dstID := common.FormatBaseTypes(DstDNI.ID)
	srcID := common.FormatBaseTypes(SrcDNI.ID)
	if common.IsSliceContainsStrNocase(SrcDNI.dstBind, dstID) &&
		common.IsSliceContainsStrNocase(DstDNI.srcBind, srcID) {
		return nil
	}

	SrcDNI_DstBd := common.AppendStr(SrcDNI.dstBind, dstID)
	DstDNI_SrcBd := common.AppendStr(DstDNI.srcBind, srcID)

	SrcDNI.DstBindList = joinInternalString(SrcDNI_DstBd)
	DstDNI.SrcBindList = joinInternalString(DstDNI_SrcBd)

	// 开始更新表
	tx := devDb.Begin()
	if err = tx.Error; err != nil {
		return err
	}

	tx.Model(&SrcDNI).Updates(&ZbDeviceNodeInfo{DstBindList: SrcDNI.DstBindList})
	if err = tx.Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Model(&DstDNI).Updates(&ZbDeviceNodeInfo{SrcBindList: DstDNI.SrcBindList})
	if err = tx.Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	if err = tx.Error; err != nil {
		tx.Rollback()
		return err
	}

	// 成功更新本地绑定
	SrcDNI.dstBind = SrcDNI_DstBd
	DstDNI.srcBind = DstDNI_SrcBd

	return nil
}

// 解绑两个设备节点
// 如果两个设备节点有两个互补集绑定的,那么将进行同时解绑
func UnZbBindDeviceNode(SrcIeeeAddr uint64, SrcNodeNum uint16,
	DstBindIeeeAddr uint64, DstBindNodeNum, BindTrunkID uint16) error {
	var SrcDNI *ZbDeviceNodeInfo
	var DstDNI *ZbDeviceNodeInfo
	var err error

	// 获取源设备节点
	if SrcDNI, err = LookupZbDeviceNodeByIN(SrcIeeeAddr, SrcNodeNum); err != nil {
		return nil
	}
	// 获取目的设备节点
	if DstDNI, err = LookupZbDeviceNodeByIN(DstBindIeeeAddr, DstBindNodeNum); err != nil {
		return nil
	}

	strBindtkid := common.FormatBaseTypes(BindTrunkID)
	// 只有源设备节点输出集和目的设备输入集互补,即都含有要绑定的集,才进行解绑定,否则认为是成功的
	if !common.IsSliceContainsStrNocase(SrcDNI.outTrunk, strBindtkid) ||
		!common.IsSliceContainsStrNocase(DstDNI.inTrunk, strBindtkid) {
		return nil
	}

	dstid := common.FormatBaseTypes(DstDNI.ID)
	srcid := common.FormatBaseTypes(SrcDNI.ID)

	// 源设备节点的<目的绑定表>不含目标设备节点
	//或 目标设备节点<源绑定表>不含源设备节点 将进行绑定解绑直接返回成功
	if !common.IsSliceContainsStrNocase(SrcDNI.dstBind, dstid) || !common.IsSliceContainsStrNocase(DstDNI.srcBind, srcid) {
		return nil
	}
	// 删除源的目标绑定 和 目标的源绑定
	SrcDNI_DstBd := common.DeleteFromSliceStr(SrcDNI.dstBind, dstid)
	DstDNI_SrcBd := common.DeleteFromSliceStr(DstDNI.srcBind, srcid)

	SrcDNI.DstBindList = joinInternalString(SrcDNI_DstBd)
	DstDNI.SrcBindList = joinInternalString(DstDNI_SrcBd)

	// 开始更新表
	tx := devDb.Begin()
	if err = tx.Error; err != nil {
		return err
	}

	tx.Model(&SrcDNI).Updates(&ZbDeviceNodeInfo{DstBindList: SrcDNI.DstBindList})
	if err = tx.Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Model(&DstDNI).Updates(&ZbDeviceNodeInfo{SrcBindList: DstDNI.SrcBindList})
	if err = tx.Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	if err = tx.Error; err != nil {
		tx.Rollback()
		return err
	}
	// 成功更新本地
	SrcDNI.dstBind = SrcDNI_DstBd
	DstDNI.srcBind = DstDNI_SrcBd

	return nil
}

// 找到绑定表的所有目标设备节点
func BindFindZbDeviceNodeByNN(NwkAddr, NodeNum, trunkID uint16) ([]*ZbDeviceNodeInfo, error) {
	src, err := LookupZbDeviceNodeByNN(NwkAddr, NodeNum)
	if err != nil {
		return nil, err
	}
	strTkid := common.FormatBaseTypes(trunkID)
	// 源设备 是否包含输出集
	if common.IsSliceContainsStrNocase(src.outTrunk, strTkid) != true {
		return nil, errors.New("不包含对应的集")
	}

	dni := make([]*ZbDeviceNodeInfo, 0, len(src.dstBind))
	for _, id := range src.dstBind {
		tmpdni, err := LookupZbDeviceNodeByID(id)
		if err != nil {
			continue
		}

		// 只有目标设备含有输入集才加入
		if common.IsSliceContainsStrNocase(tmpdni.inTrunk, strTkid) {
			dni = append(dni, tmpdni)
		}
	}

	return dni, nil
}

// 找到绑定表的所有目标设备节点
func BindFindZbDeviceNodeByIN(ieeeAddr uint64, NodeNum, trunkID uint16) ([]*ZbDeviceNodeInfo, error) {
	src, err := LookupZbDeviceNodeByIN(ieeeAddr, NodeNum)
	if err != nil {
		return nil, err
	}
	strTkid := common.FormatBaseTypes(trunkID)
	// 源设备 是否包含输出集
	if !common.IsSliceContainsStrNocase(src.outTrunk, strTkid) {
		return nil, errors.New("不包含对应的集")
	}

	dni := make([]*ZbDeviceNodeInfo, 0, len(src.dstBind))
	for _, id := range src.dstBind {
		tmpdni, err := LookupZbDeviceNodeByID(id)
		if err != nil {
			continue
		}

		// 只有目标设备含有输入集才加入
		if common.IsSliceContainsStrNocase(tmpdni.inTrunk, strTkid) {
			dni = append(dni, tmpdni)
		}
	}

	return dni, nil
}

// 解析内部的string分隔
func (this *ZbDeviceNodeInfo) parseInternalString() {
	this.inTrunk = splitInternalString(this.InTrunkList)
	this.outTrunk = splitInternalString(this.OutTrunkList)
	this.srcBind = splitInternalString(this.SrcBindList)
	this.dstBind = splitInternalString(this.DstBindList)
}

// 获取设备节点id
func (this *ZbDeviceNodeInfo) GetID() uint {
	return this.ID
}

// 获取设备节点网络地址
func (this *ZbDeviceNodeInfo) GetNwkAddr() uint16 {
	return this.NwkAddr
}

// 获取设备节点节点号
func (this *ZbDeviceNodeInfo) GetNodeNum() uint16 {
	return this.NodeNo
}

// 获取设备节点Ieee地址
func (this *ZbDeviceNodeInfo) GetIeeeAddr() uint64 {
	return this.IeeeAddr
}

// 获取设备节点 集id表
func (this *ZbDeviceNodeInfo) GetTrunkIDList() (inTrunk, outTrunk []string) {
	return this.inTrunk, this.outTrunk
}

// 获取设备节点源绑定id列表
func (this *ZbDeviceNodeInfo) GetBindList() (srcBind, dstBind []string) {
	return this.srcBind, this.dstBind
}

// 根据网络地址找到设备
func LookupZbDeviceByNwkAddr(nwkAddr uint16) (*ZbDeviceInfo, error) {
	oInfo := &ZbDeviceInfo{}
	if devDb.Where(&ZbDeviceInfo{NwkAddr: nwkAddr}).First(oInfo).RecordNotFound() {
		return nil, gorm.ErrRecordNotFound
	}

	return oInfo, nil
}

// 根据ieee地址找到设备
func LookupZbDeviceByIeeeAddr(ieeeAddr uint64) (*ZbDeviceInfo, error) {
	oInfo := &ZbDeviceInfo{}
	if devDb.Where(&ZbDeviceInfo{IeeeAddr: ieeeAddr}).First(oInfo).RecordNotFound() {
		return nil, gorm.ErrRecordNotFound
	}

	return oInfo, nil
}

// 获取ieee地址
func (this *ZbDeviceInfo) GetIeeeAddr() uint64 {
	return this.IeeeAddr
}

// 获取网络地址
func (this *ZbDeviceInfo) GetNwkAddr() uint16 {
	return this.NwkAddr
}

// 获取设备能力
func (this *ZbDeviceInfo) GetCapacity() byte {
	return this.Capacity
}

// 获取设备的id
func (this *ZbDeviceInfo) GetID() uint {
	return this.ID
}

// 获取设备的产品id
func (this *ZbDeviceInfo) GetProductID() int {
	return this.ProductId
}

func (this *ZbDeviceInfo) updateCapacity(newCapacityValue byte) error {
	if newCapacityValue == this.Capacity {
		return nil
	}

	devDb.Model(this).Update("capacity", newCapacityValue)
	return devDb.Error
}

// 更新设备和设备节点所有节点的网络地址
func (this *ZbDeviceInfo) updateZbDeviceAndNodeNwkAddr(NewnwkAddr uint16) error {
	var err error

	tx := devDb.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	// 更新设备网络地址
	if err = tx.Model(this).Updates(&ZbDeviceInfo{NwkAddr: NewnwkAddr}).Error; err != nil {
		tx.Rollback()
		return err
	}
	//更新所有节点网络地址
	if err = tx.Model(&ZbDeviceNodeInfo{}).Where(&ZbDeviceNodeInfo{IeeeAddr: this.IeeeAddr}).Updates(&ZbDeviceNodeInfo{NwkAddr: NewnwkAddr}).Error; err != nil {
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
func (this *ZbDeviceInfo) createZbDeveiceAndNode() error {
	// 查询对应产品
	pdt, err := LookupZbProduct(this.ProductId)
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
		dnode := &ZbDeviceNodeInfo{
			IeeeAddr:  this.IeeeAddr,
			NodeNo:    uint16(i + 1),
			NwkAddr:   this.NwkAddr,
			ProductId: this.ProductId,
		}

		dnode.InTrunkList = strings.Replace(
			strings.Trim(fmt.Sprint(v.InTrunk), "[]"), " ", ",", -1)
		dnode.OutTrunkList = strings.Replace(
			strings.Trim(fmt.Sprint(v.OutTrunk), "[]"), " ", ",", -1)
		dnode.SrcBindList = ""
		dnode.DstBindList = ""

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
func (this *ZbDeviceInfo) deleteZbDeveiceAndNode() error {
	tx := devDb.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var devNodes []ZbDeviceNodeInfo
	tx.Find(&devNodes, &ZbDeviceNodeInfo{IeeeAddr: this.IeeeAddr})

	for _, v := range devNodes {
		v.parseInternalString() // 将集与绑定表解析一下
		vid := common.FormatBaseTypes(v.ID)
		// 是否有源绑定,有则删除每一个 源的目标绑定
		if len(v.srcBind) > 0 {
			for _, tv := range v.srcBind { // 扫描每一个源 的目标绑定,让它删除对应id
				tmpdevNode, err := lookupZbDeviceNodeByID(tx, tv)
				if err == nil && len(tmpdevNode.dstBind) > 0 &&
					common.IsSliceContainsStrNocase(tmpdevNode.dstBind, vid) {
					tmpdevNode.dstBind = common.DeleteFromSliceStrAll(tmpdevNode.dstBind, vid)
					tmpdevNode.DstBindList = joinInternalString(tmpdevNode.dstBind)

					tx.Model(tmpdevNode).
						Updates(&ZbDeviceNodeInfo{DstBindList: tmpdevNode.DstBindList})
					if tx.Error != nil {
						tx.Rollback()
						return tx.Error
					}
				}
			}
		}

		// 是否有目标绑定, 有则删除每一个 目标的源绑定
		if len(v.dstBind) > 0 {
			for _, tv := range v.dstBind { // 扫描每一个目标 的源绑定,让它删除对应id
				tmpdevNode, err := lookupZbDeviceNodeByID(tx, tv)
				if err == nil && len(tmpdevNode.srcBind) > 0 &&
					common.IsSliceContainsStrNocase(tmpdevNode.srcBind, vid) {
					tmpdevNode.srcBind = common.DeleteFromSliceStrAll(tmpdevNode.srcBind, vid)
					tmpdevNode.SrcBindList = joinInternalString(tmpdevNode.srcBind)

					tx.Model(tmpdevNode).
						Updates(&ZbDeviceNodeInfo{SrcBindList: tmpdevNode.SrcBindList})
					if tx.Error != nil {
						tx.Rollback()
						return tx.Error
					}
				}
			}
		}

		tx.Unscoped().Delete(v)
		if tx.Error != nil {
			tx.Rollback()
			return tx.Error
		}
	}

	if err := tx.Unscoped().Delete(this).Error; err != nil {
		tx.Rollback()
		return err
	}

	if tx.Commit().Error != nil {
		tx.Rollback()
		return tx.Error
	}

	return nil
}

// 添加设备信息,包括所有设备节点,如果已存在,则更新相应参数
func UpdateZbDeviceAndNode(ieeeAddr uint64, nwkAddr uint16, capacity byte, productID int) error {
	var err error

	dev, err := LookupZbDeviceByIeeeAddr(ieeeAddr)
	if err != nil {
		return (&ZbDeviceInfo{
			IeeeAddr:  ieeeAddr,
			NwkAddr:   nwkAddr,
			Capacity:  capacity,
			ProductId: productID,
		}).createZbDeveiceAndNode()
	}

	// 已存在,看是否要更新相应参数
	if dev.ProductId != productID {
		//删除设备和设备所有的节点
		if err = dev.deleteZbDeveiceAndNode(); err != nil {
			return err
		}

		//重建所有设备和节点
		return (&ZbDeviceInfo{
			IeeeAddr:  ieeeAddr,
			NwkAddr:   nwkAddr,
			Capacity:  capacity,
			ProductId: productID,
		}).createZbDeveiceAndNode()
	}

	// 更新能力属性
	if dev.Capacity != capacity {
		if err = dev.updateCapacity(capacity); err != nil {
			return err
		}
	}

	// 更新对应网络址
	if dev.NwkAddr != nwkAddr {
		if err := dev.updateZbDeviceAndNodeNwkAddr(nwkAddr); err != nil {
			return err
		}
	}

	return nil
}

//删除这个设备及设备的所有节点
func DeleteZbDeveiceAndNode(ieeeAddr uint64) error {
	dev, err := LookupZbDeviceByIeeeAddr(ieeeAddr)
	if err != nil {
		return nil
	}

	return dev.deleteZbDeveiceAndNode()
}
