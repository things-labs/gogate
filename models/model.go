package models

import (
	"fmt"
	"os"
	"path"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/thinkgos/gogate/misc"
	"github.com/thinkgos/memlog"
	"github.com/thinkgos/utils"
)

const (
	_DbName   = "data/models.db"
	_DbDriver = "sqlite3"
)

// DbTableInitFunc 数据表初始化函数
type DbTableInitFunc func() error

var (
	db              *gorm.DB
	dbTableInitList []DbTableInitFunc
)

// DbInit 数据库初始化,注册相应模型
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
		return fmt.Errorf("%v: db(%s-%s) open failed", err, _DbDriver, _DbName)
	}
	// default disable
	db.LogMode(misc.APPConfig.OrmDbLog)
	// db.LogMode(true)

	for _, initF := range dbTableInitList {
		if err = initF(); err != nil {
			errs = err
			memlog.Error(err)
		}
	}

	return errs
}

// RegisterDbTableInitFunc 提供数据表注册初始函数
func RegisterDbTableInitFunc(f DbTableInitFunc) {
	if f != nil {
		dbTableInitList = append(dbTableInitList, f)
	}
}
