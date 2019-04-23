package models

import (
	"errors"
)

var (
	ErrInvalidParameter      = errors.New("invalid parameter")             // 无效参数
	ErrProductNotExist       = errors.New("product not exist")             // 产品不存在
	ErrDeviceNotExist        = errors.New("device not exist")              // 设备不存在
	ErrNotContainTrunk       = errors.New("not contain the trunk id")      // 没有包含指定集id
	ErrTrunkNotComplementary = errors.New("trunkI id Not a complementary") // 集id不是互补的
)
