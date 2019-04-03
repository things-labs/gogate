package models

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/thinkgos/utils"
)

type User struct {
	gorm.Model
	Uid int64 `gorm:"UNIQUE;NOT NULL"`
}
type localUser struct {
	sync.RWMutex
	tab []int64
}

var lUser *localUser

func init() {
	db.AutoMigrate(&User{})
	lUser = new(localUser)
	users := getUsers()
	lUser.tab = make([]int64, 0, len(users))
	for _, v := range users {
		lUser.tab = append(lUser.tab, v.Uid)
	}
}

func HasUser(uid int64) bool {
	if uid == 0 {
		return true
	}
	lUser.RLock()
	b := utils.IsSliceContainsInt64(lUser.tab, uid)
	lUser.RUnlock()
	return b
}

func AddUser(uid int64) error {
	if (uid == 0) || HasUser(uid) {
		return nil
	}

	if err := db.Create(&User{Uid: uid}).Error; err != nil {
		return err
	}
	lUser.Lock()
	lUser.tab = append(lUser.tab, uid)
	lUser.Unlock()

	return nil
}

func DeleteUser(uid int64) error {
	if (uid == 0) || !HasUser(uid) {
		return nil
	}
	user := &User{Uid: uid}
	if err := db.Where(user).Unscoped().Delete(user).Error; err != nil {
		return err
	}
	lUser.Lock()
	lUser.tab = utils.DeleteFromSliceInt64All(lUser.tab, uid)
	lUser.Unlock()
	return nil
}

func getUsers() []User {
	users := []User{}
	db.Find(&users)
	return users
}

func GetUsers() []int64 {
	lUser.RLock()
	tb := make([]int64, len(lUser.tab))
	copy(tb, lUser.tab)
	lUser.RUnlock()
	return tb
}
