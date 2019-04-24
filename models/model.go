package models

import (
	"os"
	"path"

	"github.com/thinkgos/gogate/misc"
	"github.com/thinkgos/utils"

	"github.com/astaxie/beego/logs"
	"github.com/go-ini/ini"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

const (
	_DB_NAME   = "data/models.db"
	_DB_DRIVER = "sqlite3"
)

type DbTableInitFunc func() error

var (
	db              *gorm.DB
	dbTableInitList []DbTableInitFunc
)

// 数据库初始化,注册相应模型
func DbInit() error {
	var err error
	var errs error

	// 判断目录是否存在,不存在着创建对应的所有目录
	if !utils.IsExist(_DB_NAME) {
		if err = os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm); err != nil {
			return err
		}
		if _, err = os.Create(_DB_NAME); err != nil {
			return err
		}
	}

	if db, err = gorm.Open(_DB_DRIVER, _DB_NAME); err != nil {
		return errors.Wrapf(err, "db(%s-%s) open failed", _DB_DRIVER, _DB_NAME)
	}
	//default disable
	db.LogMode(misc.APPCfg.Section(ini.DefaultSection).Key("ormDbLog").MustBool(false))
	//db.LogMode(true)

	for _, initF := range dbTableInitList {
		if err = initF(); err != nil {
			errs = err
			logs.Error(err)
		}
	}

	return errs
}

// 提供数据表注册初始函数
func RegisterDbTableInitFunction(f DbTableInitFunc) {
	if f != nil {
		dbTableInitList = append(dbTableInitList, f)
	}
}
