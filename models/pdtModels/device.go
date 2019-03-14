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

type DeviceInfo struct {
	ID        uint
	Sn        string
	ProductId int
}

const (
	deviceInfo_sql = `CREATE TABLE "device_infos" (
			"id" integer primary key autoincrement,
			"sn" bigint NOT NULL,
			"product_id" integer NOT NULL,
			UNIQUE(sn,product_id) ON CONFLICT FAIL)`
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
		devDb.Raw(deviceInfo_sql).Scan(&DeviceInfo{})
	}
}

// 是否对应产品的设备
func HasDevice(sn string, pid int) bool {
	_, ok := LookupProduct(pid)
	if !ok {
		return false
	}

	// switch p.Types {
	// case ProductTypes_Zigbee:
	// 	return hasZbDevice(sn, pid)
	// }

	return hasGeneralDevice(sn, pid)
}

// 创建对应产品的设备
func CreateDevice(sn string, pid int) error {
	_, ok := LookupProduct(pid)
	if !ok {
		return ErrProductNotExist
	}

	// switch p.Types {
	// case ProductTypes_Zigbee:
	// 	return createZbDevice(sn, pid)
	// }

	return createGeneralDevice(sn, pid)
}

// 删除对应产品的设备
func DeleteDevice(sn string, pid int) error {
	_, ok := LookupProduct(pid)
	if !ok {
		return ErrProductNotExist
	}
	return deleteGeneralDevice(sn, pid)
}

// 是否有通用对应的设备
func hasGeneralDevice(sn string, pid int) bool {
	if devDb.Where(&DeviceInfo{Sn: sn, ProductId: pid}).First(&DeviceInfo{}).RecordNotFound() {
		return false
	}

	return true
}

// 创建通用设备
func (this *DeviceInfo) createGeneralDevice() error {
	return devDb.Create(this).Error
}

// 创建通用设备
func createGeneralDevice(sn string, pid int) error {
	return (&DeviceInfo{Sn: sn, ProductId: pid}).createGeneralDevice()
}

// 删除通用设备
func (this *DeviceInfo) deleteGeneralDevice() error {
	return devDb.Unscoped().Delete(this).Error
}

// 删除通用设备
func deleteGeneralDevice(sn string, pid int) error {
	return (&DeviceInfo{Sn: sn, ProductId: pid}).deleteGeneralDevice()
}
