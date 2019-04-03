package models

import (
	"os"
	"path"

	"github.com/Unknwon/com"
	"github.com/jinzhu/gorm"
)

const (
	_DB_NAME   = "data/models.db"
	_DB_DRIVER = "sqlite3"
)

var db *gorm.DB

func init() {
	var err error

	// 判断目录是否存在,不存在着创建对应的所有目录
	if !com.IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}

	if db, err = gorm.Open(_DB_DRIVER, _DB_NAME); err != nil {
		panic("devmodels: gorm open failed," + err.Error())
	}
	//default disable
	//devDb.LogMode(misc.APPCfg.MustBool(goconfig.DEFAULT_SECTION, "ormDbLog", false))
	db.LogMode(true)

	db.AutoMigrate(&ZbDeviceInfo{})
	if db.Error != nil {
		panic("devmodels: gorm AutoMigrate failed," + err.Error())
	}

	if !db.HasTable("general_device_infos") {
		db.Raw(generaldeviceInfo_sql).Scan(&GeneralDeviceInfo{})
	}
	if !db.HasTable("zb_device_node_infos") {
		db.Raw(zbDeviceNodeInfos_Sql).Scan(&ZbDeviceNodeInfo{})
	}
}
