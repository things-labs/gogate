package devmodels

import (
	"os"
	"path"

	"github.com/Unknwon/com"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

const (
	_DB_NAME   = "data/devmodels.db"
	_DB_DRIVER = "sqlite3"
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

	if devDb, err = gorm.Open(_DB_DRIVER, _DB_NAME); err != nil {
		panic("devmodels: gorm open failed," + err.Error())
	}
	//default disable
	//devDb.LogMode(misc.APPCfg.MustBool(goconfig.DEFAULT_SECTION, "ormDbLog", false))
	devDb.LogMode(true)

	devDb.AutoMigrate(&ZbDeviceInfo{})
	if devDb.Error != nil {
		panic("devmodels: gorm AutoMigrate failed," + err.Error())
	}

	if !devDb.HasTable("general_device_infos") {
		devDb.Raw(generaldeviceInfo_sql).Scan(&GeneralDeviceInfo{})
	}
	if !devDb.HasTable("zb_device_node_infos") {
		devDb.Raw(zbDeviceNodeInfos_Sql).Scan(&ZbDeviceNodeInfo{})
	}
}

// 是否有通用对应的设备
func HasGeneralDevice(pid int, sn string) bool {
	_, err := LookupProduct(pid)
	if err != nil {
		return false
	}

	if devDb.Where(&GeneralDeviceInfo{ProductId: pid, Sn: sn}).
		First(&GeneralDeviceInfo{}).RecordNotFound() {
		return false
	}
	return true
}

// 创建通用设备
func CreateGeneralDevice(pid int, sn string) error {
	_, err := LookupProduct(pid)
	if err != nil {
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
	_, err := LookupProduct(pid)
	if err != nil {
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
	_, err := LookupProduct(pid)
	if err != nil {
		return []GeneralDeviceInfo{}
	}
	devs := []GeneralDeviceInfo{}
	devDb.Where(&GeneralDeviceInfo{ProductId: pid}).Find(&devs)
	return devs
}