package devmodels

import (
	"errors"
)

var (
	ErrProductNotExist = errors.New("prouduct not exist")
	ErrDeviceNotExist  = errors.New("device not exist")
	ErrNotContainTrunk = errors.New("not contain the trunk id")
)
