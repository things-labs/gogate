package models

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/thinkgos/utils"
)

const (
	SupperUser = 0 // 超级用户,超级权限,必定存在
)

type User struct {
	gorm.Model
	Uid int64 `gorm:"UNIQUE;NOT NULL"`
}

type UserInfo struct {
	sync.RWMutex
	tab []int64
}

var localUser *UserInfo

func init() {
	RegisterDbTableInitFunction(func() error {
		if err := db.AutoMigrate(&User{}).Error; err != nil {
			return err
		}
		localUser = new(UserInfo)
		users := []User{}
		db.Find(&users)
		localUser.tab = make([]int64, 0, len(users))
		for _, v := range users {
			localUser.tab = append(localUser.tab, v.Uid)
		}
		return nil
	})
}

// 是否有对应用户,用户0为超级用户,永远存在
func HasUser(uid int64) bool {
	if uid == SupperUser {
		return true
	}
	localUser.RLock()
	b := utils.IsSliceContainsInt64(localUser.tab, uid)
	localUser.RUnlock()
	return b
}

// 添加用户
func AddUser(uid int64) error {
	if HasUser(uid) {
		return nil
	}

	if err := db.Create(&User{Uid: uid}).Error; err != nil {
		return err
	}
	localUser.Lock()
	localUser.tab = append(localUser.tab, uid)
	localUser.Unlock()

	return nil
}

// 删除用户
func DeleteUser(uid int64) error {
	if !HasUser(uid) {
		return nil
	}
	user := &User{Uid: uid}
	if err := db.Where(user).Unscoped().Delete(user).Error; err != nil {
		return err
	}
	localUser.Lock()
	localUser.tab = utils.DeleteFromSliceInt64All(localUser.tab, uid)
	localUser.Unlock()
	return nil
}

// 获取用户列表
func GetUsers() []int64 {
	tb := make([]int64, len(localUser.tab))
	localUser.RLock()
	copy(tb, localUser.tab) // 拷贝
	localUser.RUnlock()
	return tb
}
