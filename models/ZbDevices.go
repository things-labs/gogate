package models

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/thinkgos/utils"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// 设备表
type ZbDeviceInfo struct {
	gorm.Model
	Sn        string `gorm:"UNIQUE;NOT NULL"`
	NwkAddr   uint16 `gorm:"NOT NULL"`
	Capacity  byte   `gorm:"default:2"`
	ProductId int
}

// NOTE: 表内数组全部以逗号","分隔
// 节点表
type ZbDeviceNodeInfo struct {
	ID           uint   //`gorm:"primary_key"`
	Sn           string //`gorm:"NOT NULL"`
	NodeNo       byte   //`gorm:"UNIQUE;NOT NULL"`
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
	"sn" varchar(255) NOT NULL,
	"node_no" integer NOT NULL,
	"nwk_addr" integer NOT NULL,
	"product_id" integer,
	"in_trunk_list" varchar(255),
	"out_trunk_list" varchar(255),
	"src_bind_list" varchar(255),
	"dst_bind_list" varchar(255),
	"status" integer default(0),
	UNIQUE(sn,node_no) ON CONFLICT FAIL)`

// zigbee设备表初始化
func ZbDeviceDbTableInit() error {
	if err := db.AutoMigrate(&ZbDeviceInfo{}).Error; err != nil {
		return errors.Wrap(err, "db AutoMigrate failed")
	}

	if !db.HasTable("zb_device_node_infos") {
		//TODO: check error?
		db.Raw(zbDeviceNodeInfos_Sql).Scan(&ZbDeviceNodeInfo{})
	}
	return nil
}

// 采用","分隔字符串
func splitInternalString(s string) []string {
	if s == "" {
		return []string{}
	}
	// NOTE: 如果切割空字符串会返回一个空字符的数组,长度为1,这是个坑
	return strings.Split(s, ",")
}

// 采用","拼接字符串
func joinInternalString(s []string) string {
	return strings.Join(s, ",")
}

// 是否有指定的设备
func HasZbDevice(sn string, pid int) bool {
	if !HasZbProduct(pid) || len(sn) == 0 {
		return false
	}
	return db.Where(&ZbDeviceInfo{Sn: sn}).First(&ZbDeviceInfo{}).Error == nil
}

// 有指定节点(nwkAddress nodeNum)
func HasZbDeviceNode(nwkAddr uint16, nodeNum byte) bool {
	return db.Where(&ZbDeviceNodeInfo{
		NwkAddr: nwkAddr,
		NodeNo:  nodeNum,
	}).Error == nil
}

// 根据nwkAddr,nodeNum找到设备节点
func LookupZbDeviceNodeByNN(nwkAddr uint16, nodeNum byte) (*ZbDeviceNodeInfo, error) {
	o := &ZbDeviceNodeInfo{}
	if err := db.Where(&ZbDeviceNodeInfo{
		NwkAddr: nwkAddr,
		NodeNo:  nodeNum}).First(o).Error; err != nil {
		return nil, err
	}
	return o.parseInternalString(), nil
}

// 根据sn,nodeNum找到设备节点
func LookupZbDeviceNodeByIN(sn string, nodeNum byte) (*ZbDeviceNodeInfo, error) {
	if len(sn) == 0 {
		return nil, ErrInvalidParameter
	}
	o := &ZbDeviceNodeInfo{}
	if err := db.Where(&ZbDeviceNodeInfo{
		Sn:     sn,
		NodeNo: nodeNum}).First(o).Error; err != nil {
		return nil, err
	}
	return o.parseInternalString(), nil
}

// 根据id找到设备节点
func LookupZbDeviceNodeByID(id string) (*ZbDeviceNodeInfo, error) {
	return lookupZbDeviceNodeByID(db, id)
}

// 根据id找到设备节点
func lookupZbDeviceNodeByID(db *gorm.DB, id string) (*ZbDeviceNodeInfo, error) {
	idval, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		return nil, err
	}
	o := &ZbDeviceNodeInfo{}
	if err := db.First(o, idval).Error; err != nil {
		return nil, err
	}

	return o.parseInternalString(), nil
}

// 绑定两个设备 要求更新 源设备节点的<目的绑定表>和目标设备节点的<源绑定表>
func BindZbDeviceNode(SrcSn string, SrcNodeNum byte,
	DstBindSn string, DstBindNodeNum byte,
	BindTrunkID uint16) error {
	var SrcDNI *ZbDeviceNodeInfo
	var DstBindDNI *ZbDeviceNodeInfo
	var err error

	// 获取源设备节点
	if SrcDNI, err = LookupZbDeviceNodeByIN(SrcSn, SrcNodeNum); err != nil {
		return err
	}
	// 获取目的设备节点
	if DstBindDNI, err = LookupZbDeviceNodeByIN(DstBindSn, DstBindNodeNum); err != nil {
		return err
	}

	strBindTkid := utils.FormatBaseTypes(BindTrunkID)
	// 只有源设备节点输出集和目的设备输入集互补,即都含有要绑定的集,才进行绑定
	if !utils.IsSliceContainsStrNocase(SrcDNI.outTrunk, strBindTkid) ||
		!utils.IsSliceContainsStrNocase(DstBindDNI.inTrunk, strBindTkid) {
		return ErrTrunkNotComplementary
	}

	srcID := utils.FormatBaseTypes(SrcDNI.ID)
	dstBindID := utils.FormatBaseTypes(DstBindDNI.ID)
	// 源设备节点 目的绑定表 不含目标设备节点 或 目标设备节点 源绑定表 不含源设备节点 将进行绑定添加,
	// 都有直接返回成功
	if utils.IsSliceContainsStrNocase(SrcDNI.dstBind, dstBindID) &&
		utils.IsSliceContainsStrNocase(DstBindDNI.srcBind, srcID) {
		return nil
	}
	// 将目的设备id号添加到 源设备的 目的绑定表
	SrcDNI_DstBd := utils.AppendStr(SrcDNI.dstBind, dstBindID)
	// 将源设备id号添加到 目的设备的 源绑定表
	DstBindDNI_SrcBd := utils.AppendStr(DstBindDNI.srcBind, srcID)

	SrcDNI.DstBindList = joinInternalString(SrcDNI_DstBd)
	DstBindDNI.SrcBindList = joinInternalString(DstBindDNI_SrcBd)

	// 开始更新表
	tx := db.Begin()
	if err = tx.Error; err != nil {
		return err
	}

	if err = tx.Model(&SrcDNI).
		Updates(&ZbDeviceNodeInfo{DstBindList: SrcDNI.DstBindList}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Model(&DstBindDNI).
		Updates(&ZbDeviceNodeInfo{SrcBindList: DstBindDNI.SrcBindList}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	// 成功更新本地绑定
	SrcDNI.dstBind = SrcDNI_DstBd
	DstBindDNI.srcBind = DstBindDNI_SrcBd

	return nil
}

// 解绑两个设备节点
// 如果两个设备节点有两个互补集绑定的,那么将进行同时解绑
func UnZbBindDeviceNode(SrcSn string, SrcNodeNum byte,
	DstBindSn string, DstBindNodeNum byte, BindTrunkID uint16) error {
	var SrcDNI *ZbDeviceNodeInfo
	var DstBindDNI *ZbDeviceNodeInfo
	var err error

	// 获取源设备节点
	if SrcDNI, err = LookupZbDeviceNodeByIN(SrcSn, SrcNodeNum); err != nil {
		return nil
	}
	// 获取目的设备节点
	if DstBindDNI, err = LookupZbDeviceNodeByIN(DstBindSn, DstBindNodeNum); err != nil {
		return nil
	}

	strBindtkid := utils.FormatBaseTypes(BindTrunkID)
	// 只有源设备节点输出集和目的设备输入集互补,即都含有要绑定的集,才进行解绑定,否则认为是成功的
	if !utils.IsSliceContainsStrNocase(SrcDNI.outTrunk, strBindtkid) ||
		!utils.IsSliceContainsStrNocase(DstBindDNI.inTrunk, strBindtkid) {
		return nil
	}

	dstid := utils.FormatBaseTypes(DstBindDNI.ID)
	srcid := utils.FormatBaseTypes(SrcDNI.ID)

	// 源设备节点的<目的绑定表>不含目标设备节点
	//或 目标设备节点<源绑定表>不含源设备节点 将进行绑定解绑直接返回成功
	if !utils.IsSliceContainsStrNocase(SrcDNI.dstBind, dstid) ||
		!utils.IsSliceContainsStrNocase(DstBindDNI.srcBind, srcid) {
		return nil
	}
	// 删除源的目标绑定 和 目标的源绑定
	SrcDNI_DstBd := utils.DeleteFromSliceStr(SrcDNI.dstBind, dstid)
	DstBindDNI_SrcBd := utils.DeleteFromSliceStr(DstBindDNI.srcBind, srcid)

	SrcDNI.DstBindList = joinInternalString(SrcDNI_DstBd)
	DstBindDNI.SrcBindList = joinInternalString(DstBindDNI_SrcBd)

	// 开始更新表
	tx := db.Begin()
	if err = tx.Error; err != nil {
		return err
	}

	if err = tx.Model(&SrcDNI).
		Updates(&ZbDeviceNodeInfo{DstBindList: SrcDNI.DstBindList}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Model(&DstBindDNI).
		Updates(&ZbDeviceNodeInfo{SrcBindList: DstBindDNI.SrcBindList}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	// 成功更新本地
	SrcDNI.dstBind = SrcDNI_DstBd
	DstBindDNI.srcBind = DstBindDNI_SrcBd

	return nil
}

// 找到绑定表的所有目标设备节点
func BindFindZbDeviceNodeByNN(NwkAddr uint16, NodeNum byte, trunkID uint16) ([]*ZbDeviceNodeInfo, error) {
	srcDNI, err := LookupZbDeviceNodeByNN(NwkAddr, NodeNum)
	if err != nil {
		return nil, err
	}
	strTkid := utils.FormatBaseTypes(trunkID)
	// 源设备 是否包含输出集
	if !utils.IsSliceContainsStrNocase(srcDNI.outTrunk, strTkid) {
		return nil, ErrNotContainTrunk
	}

	dni := make([]*ZbDeviceNodeInfo, 0, len(srcDNI.dstBind))
	for _, id := range srcDNI.dstBind {
		tmpdni, err := LookupZbDeviceNodeByID(id)
		if err != nil {
			continue
		}

		// 只有目标设备含有输入集才加入
		if utils.IsSliceContainsStrNocase(tmpdni.inTrunk, strTkid) {
			dni = append(dni, tmpdni)
		}
	}

	return dni, nil
}

// 找到绑定表的所有目标设备节点
func BindFindZbDeviceNodeByIN(sn string, NodeNum byte, trunkID uint16) ([]*ZbDeviceNodeInfo, error) {
	src, err := LookupZbDeviceNodeByIN(sn, NodeNum)
	if err != nil {
		return nil, err
	}
	strTkid := utils.FormatBaseTypes(trunkID)
	// 源设备 是否包含输出集
	if !utils.IsSliceContainsStrNocase(src.outTrunk, strTkid) {
		return nil, ErrNotContainTrunk
	}

	dni := make([]*ZbDeviceNodeInfo, 0, len(src.dstBind))
	for _, id := range src.dstBind {
		tmpdni, err := LookupZbDeviceNodeByID(id)
		if err != nil {
			continue
		}

		// 只有目标设备含有输入集才加入
		if utils.IsSliceContainsStrNocase(tmpdni.inTrunk, strTkid) {
			dni = append(dni, tmpdni)
		}
	}

	return dni, nil
}

// 解析内部的string分隔
func (this *ZbDeviceNodeInfo) parseInternalString() *ZbDeviceNodeInfo {
	this.inTrunk = splitInternalString(this.InTrunkList)
	this.outTrunk = splitInternalString(this.OutTrunkList)
	this.srcBind = splitInternalString(this.SrcBindList)
	this.dstBind = splitInternalString(this.DstBindList)
	return this
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
func (this *ZbDeviceNodeInfo) GetNodeNum() byte {
	return this.NodeNo
}

// 获取设备节点sn(Ieee地址)
func (this *ZbDeviceNodeInfo) GetSn() string {
	return this.Sn
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
	err := db.Where(&ZbDeviceInfo{NwkAddr: nwkAddr}).First(oInfo).Error
	return oInfo, err
}

// 根据sn找到设备
func LookupZbDeviceByIeeeAddr(sn string) (*ZbDeviceInfo, error) {
	if len(sn) == 0 {
		return nil, ErrInvalidParameter
	}
	oInfo := &ZbDeviceInfo{}
	err := db.Where(&ZbDeviceInfo{Sn: sn}).First(oInfo).Error
	return oInfo, err
}

// 获取sn
func (this *ZbDeviceInfo) GetSn() string {
	return this.Sn
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

func (this *ZbDeviceInfo) updateCapacity(newVal byte) error {
	if newVal == this.Capacity {
		return nil
	}
	return db.Model(this).Update("capacity", newVal).Error
}

// 更新设备和设备节点所有节点的网络地址
func (this *ZbDeviceInfo) updateZbDeviceAndNodeNwkAddr(NewnwkAddr uint16) error {
	var err error

	tx := db.Begin()
	if err = tx.Error; err != nil {
		return err
	}
	// 更新设备网络地址
	if err = tx.Model(this).
		Updates(&ZbDeviceInfo{NwkAddr: NewnwkAddr}).Error; err != nil {
		tx.Rollback()
		return err
	}
	//更新所有节点网络地址
	if err = tx.Model(&ZbDeviceNodeInfo{}).
		Where(&ZbDeviceNodeInfo{Sn: this.Sn}).
		Updates(&ZbDeviceNodeInfo{NwkAddr: NewnwkAddr}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return err
}

// 创建设备和设备所有的节点,失败将不建立
func (this *ZbDeviceInfo) createZbDeveiceAndNode() error {
	// 查询对应产品
	pdt, err := LookupZbProduct(this.ProductId)
	if err != nil {
		return err
	}

	devNode := pdt.GetDeviceNodeDscList()

	tx := db.Begin()
	if err = tx.Error; err != nil {
		return err
	}
	// 创建设备

	if err = tx.Create(this).Error; err != nil {
		tx.Rollback()
		return err
	}

	//创建除保留节点(0)外的所有节点
	for i, v := range devNode {
		dnode := &ZbDeviceNodeInfo{
			Sn:        this.Sn,
			NodeNo:    byte(i + 1),
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
	}

	return nil
}

// 根据ieee地址删除设备,包含所有的设备节点
func (this *ZbDeviceInfo) DeleteZbDeveiceAndNode() error {
	var err error

	tx := db.Begin()
	if err = tx.Error; err != nil {
		return err
	}

	var devNodes []ZbDeviceNodeInfo
	tx.Find(&devNodes, &ZbDeviceNodeInfo{Sn: this.Sn})

	for _, v := range devNodes {
		v.parseInternalString() // 将集与绑定表解析一下
		vid := utils.FormatBaseTypes(v.ID)
		// 是否有源绑定,有则删除每一个 源的目标绑定
		if len(v.srcBind) > 0 {
			for _, tv := range v.srcBind { // 扫描每一个源 的目标绑定,让它删除对应id
				tmpdevNode, err := lookupZbDeviceNodeByID(tx, tv)
				if err == nil && len(tmpdevNode.dstBind) > 0 &&
					utils.IsSliceContainsStrNocase(tmpdevNode.dstBind, vid) {
					tmpdevNode.dstBind = utils.DeleteFromSliceStrAll(tmpdevNode.dstBind, vid)
					tmpdevNode.DstBindList = joinInternalString(tmpdevNode.dstBind)

					err = tx.Model(tmpdevNode).
						Updates(&ZbDeviceNodeInfo{DstBindList: tmpdevNode.DstBindList}).Error
					if err != nil {
						tx.Rollback()
						return err
					}
				}
			}
		}

		// 是否有目标绑定, 有则删除每一个 目标的源绑定
		if len(v.dstBind) > 0 {
			for _, tv := range v.dstBind { // 扫描每一个目标 的源绑定,让它删除对应id
				tmpdevNode, err := lookupZbDeviceNodeByID(tx, tv)
				if err == nil && len(tmpdevNode.srcBind) > 0 &&
					utils.IsSliceContainsStrNocase(tmpdevNode.srcBind, vid) {
					tmpdevNode.srcBind = utils.DeleteFromSliceStrAll(tmpdevNode.srcBind, vid)
					tmpdevNode.SrcBindList = joinInternalString(tmpdevNode.srcBind)

					err = tx.Model(tmpdevNode).
						Updates(&ZbDeviceNodeInfo{SrcBindList: tmpdevNode.SrcBindList}).Error
					if err != nil {
						tx.Rollback()
						return err
					}
				}
			}
		}

		if err = tx.Unscoped().Delete(v).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err = tx.Unscoped().Delete(this).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// 添加设备信息,包括所有设备节点,如果已存在,则更新相应参数
func UpdateZbDeviceAndNode(sn string, nwkAddr uint16, capacity byte, productID int) error {
	var err error

	if !HasZbProduct(productID) || len(sn) == 0 {
		return ErrProductNotExist
	}
	dev, err := LookupZbDeviceByIeeeAddr(sn)
	if err != nil {
		return (&ZbDeviceInfo{
			Sn:        sn,
			NwkAddr:   nwkAddr,
			Capacity:  capacity,
			ProductId: productID,
		}).createZbDeveiceAndNode()
	}

	// 已存在,看是否要更新相应参数
	if dev.ProductId != productID {
		//删除设备和设备所有的节点
		if err = dev.DeleteZbDeveiceAndNode(); err != nil {
			return err
		}

		//重建所有设备和节点
		return (&ZbDeviceInfo{
			Sn:        sn,
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
func DeleteZbDeveiceAndNode(sn string) error {
	dev, err := LookupZbDeviceByIeeeAddr(sn)
	if err != nil {
		return nil
	}

	return dev.DeleteZbDeveiceAndNode()
}

// 地址转换成字符串,全大写
func ToHexString(v uint64) string {
	return strings.ToUpper(hex.EncodeToString(utils.Big_Endian.Putuint64(v)))
}
