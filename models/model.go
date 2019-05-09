package models

import (
	"os"
	"path"

	"github.com/thinkgos/utils"
	ini "gopkg.in/ini.v1"

	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/thinkgos/gogate/misc"

	_ "github.com/mattn/go-sqlite3"
)

const (
	_DbName   = "data/models.db"
	_DbDriver = "sqlite3"
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
	if !utils.IsExist(_DbName) {
		if err = os.MkdirAll(path.Dir(_DbName), os.ModePerm); err != nil {
			return err
		}
		if _, err = os.Create(_DbName); err != nil {
			return err
		}
	}

	if db, err = gorm.Open(_DbDriver, _DbName); err != nil {
		return errors.Wrapf(err, "db(%s-%s) open failed", _DbDriver, _DbName)
	}
	//default disable
	db.LogMode(misc.APPCfg.Section(ini.DEFAULT_SECTION).Key("ormDbLog").MustBool(false))
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
