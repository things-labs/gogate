package devmodels

import (
	"errors"
)

var (
	ErrProductNotExist       = errors.New("prouduct not exist")            // 产品不存在
	ErrDeviceNotExist        = errors.New("device not exist")              // 设备不存在
	ErrNotContainTrunk       = errors.New("not contain the trunk id")      // 没有包含指定集id
	ErrTrunkNotComplementary = errors.New("trunkI id Not a complementary") // 集id不是互补的
)
