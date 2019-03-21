package devmodels

import (
	"errors"
)

var (
	ErrProductNotExist = errors.New("prouduct not exist")
	ErrDeviceNotExist  = errors.New("device not exist")
)
