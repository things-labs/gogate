package pdtModels

import (
	"os"
	"path"

	"github.com/Unknwon/com"
	"github.com/jinzhu/gorm"
)

const (
	_DB_NAME  = "data/pdtModels.db"
	_DB_DRIVE = "sqlite3"
)

type GeneralDeviceInfo struct {
	ID        uint
	ProductId int
	Sn        string
}

const (
	generaldeviceInfo_sql = `CREATE TABLE "general_device_infos" (
			"id" integer primary key autoincrement,
			"product_id" integer NOT NULL,
			"sn" bigint NOT NULL,
			UNIQUE(product_id,sn) ON CONFLICT FAIL)`
)

var devDb *gorm.DB

func init() {
	var err error

	// 判断目录是否存在,不存在着创建对应的所有目录
	if !com.IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}

	if devDb, err = gorm.Open(_DB_DRIVE, _DB_NAME); err != nil {
		panic("pdtModels: gorm open failed," + err.Error())
	}
	//default disable
	//devDb.LogMode(misc.APPCfg.MustBool(goconfig.DEFAULT_SECTION, "ormDbLog", false))
	//devDb.LogMode(true)

	devDb.AutoMigrate(&ZbDeviceInfo{})
	if devDb.Error != nil {
		panic("pdtModels: gorm AutoMigrate failed," + err.Error())
	}

	if !devDb.HasTable("device_node_infos") {
		devDb.Raw(zbDeviceNodeInfos_Sql).Scan(&ZbDeviceNodeInfo{})
	}
	if !devDb.HasTable("device_infos") {
		devDb.Raw(generaldeviceInfo_sql).Scan(&GeneralDeviceInfo{})
	}
}

// 是否有通用对应的设备
func HasGeneralDevice(pid int, sn string) bool {
	_, ok := LookupProduct(pid)
	if !ok {
		return false
	}

	if devDb.Where(&GeneralDeviceInfo{ProductId: pid, Sn: sn}).First(&GeneralDeviceInfo{}).RecordNotFound() {
		return false
	}

	return true
}

// 创建通用设备
func CreateGeneralDevice(pid int, sn string) error {
	_, ok := LookupProduct(pid)
	if !ok {
		return ErrProductNotExist
	}

	return (&GeneralDeviceInfo{ProductId: pid, Sn: sn}).CreateGeneralDevice()
}

// 创建通用设备
func (this *GeneralDeviceInfo) CreateGeneralDevice() error {
	return devDb.Create(this).Error
}

// 删除通用设备
func DeleteGeneralDevice(pid int, sn string) error {
	_, ok := LookupProduct(pid)
	if !ok {
		return ErrProductNotExist
	}
	return (&GeneralDeviceInfo{ProductId: pid, Sn: sn}).DeleteGeneralDevice()
}

// 删除通用设备
func (this *GeneralDeviceInfo) DeleteGeneralDevice() error {
	return devDb.Where(this).Unscoped().Delete(this).Error
}

// 查找通用设备
func FindGeneralDevice(pid int) []GeneralDeviceInfo {
	_, ok := LookupProduct(pid)
	if !ok {
		return []GeneralDeviceInfo{}
	}

	devs := []GeneralDeviceInfo{}
	devDb.Where(&GeneralDeviceInfo{ProductId: pid}).Find(devs)
	return devs
}
