package models

import (
	"slices"
	"sync"

	"github.com/jinzhu/gorm"
)

// 超级用户
const (
	SupperUser = "0" // 超级用户,超级权限,必定存在
)

// User 用户表
type User struct {
	gorm.Model
	UID string `gorm:"UNIQUE;NOT NULL"`
}

// UserInfo 用户信息
type UserInfo struct {
	sync.RWMutex
	tab []string
}

var localUser *UserInfo

// UserDbTableInit 用户表初始化
func UserDbTableInit() error {
	if err := db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	localUser = new(UserInfo)
	users := []User{}
	db.Find(&users)
	localUser.tab = make([]string, 0, len(users))
	for _, v := range users {
		localUser.tab = append(localUser.tab, v.UID)
	}
	return nil
}

// HasUser 是否有对应用户,用户0为超级用户,永远存在
func HasUser(uid string) bool {
	if uid == SupperUser {
		return true
	}
	localUser.RLock()
	b := slices.Contains(localUser.tab, uid)
	localUser.RUnlock()
	return b
}

// AddUser 添加用户
func AddUser(uid string) error {
	if HasUser(uid) {
		return nil
	}

	if err := db.Create(&User{UID: uid}).Error; err != nil {
		return err
	}
	localUser.Lock()
	localUser.tab = append(localUser.tab, uid)
	localUser.Unlock()

	return nil
}

// DeleteUser 删除用户
func DeleteUser(uid string) error {
	if !HasUser(uid) {
		return nil
	}
	user := &User{UID: uid}
	if err := db.Where(user).Unscoped().Delete(user).Error; err != nil {
		return err
	}
	localUser.Lock()
	localUser.tab = slices.DeleteFunc(localUser.tab, func(s string) bool {
		return s == uid
	})
	localUser.Unlock()
	return nil
}

// GetUsers 获取用户列表
func GetUsers() []string {
	tb := make([]string, len(localUser.tab))
	localUser.RLock()
	copy(tb, localUser.tab) // 拷贝
	localUser.RUnlock()
	return tb
}
