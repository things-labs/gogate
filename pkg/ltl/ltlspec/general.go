package ltlspec

import (
	"fmt"
)

const (
	/*** Basic Attribute***/
	ATTRID_BASIC_LTL_VERSION = iota
	ATTRID_BASIC_APPL_VERSION
	ATTRID_BASIC_HW_VERSION
	ATTRID_BASIC_MANUFACTURER_NAME
	ATTRID_BASIC_BUILDDATE_CODE
	ATTRID_BASIC_PRODUCT_ID
	ATTRID_BASIC_SERIAL_NUMBER
	ATTRID_BASIC_POWER_SOURCE
)

const (
	// for basic power source
	POWERSOURCE_UNKOWN = iota
	POWERSOURCE_SINGLE_PHASE
	POWERSOURCE_THREE_PHASE
	POWERSOURCE_DC
	POWERSOURCE_BATTERY
	POWERSOURCE_EMERGENCY
	// Bit b7 indicates whether the device has a secondary power source in the form of a battery backup

	/*** Power Source Attribute bits  ***/
	POWER_SOURCE_PRIMARY   = 0x7F
	POWER_SOURCE_SECONDARY = 0x80
	/*** Basic Trunck Commands ***/
	COMMAND_BASIC_RESET_FACT_DEFAULT = 0x00
	COMMAND_BASIC_REBOOT_DEVICE      = 0x01
	COMMAND_BASIC_IDENTIFY           = 0x02
)
const (
	/*** On/Off Switch***/
	ATTRID_ONOFF_STATUS = 0x0000
	/*** On/Off Switch Trunck Commands ***/
	COMMAND_ONOFF_OFF    = 0x00
	COMMAND_ONOFF_ON     = 0x01
	COMMAND_ONOFF_TOGGLE = 0x02
)
const (
	/*** Level Control  ***/
	ATTRID_LEVEL_CURRENT_LEVEL = 0x0000
	/*** Level Control Commands ***/
	COMMAND_LEVEL_MOVE_TO_LEVEL = iota
	COMMAND_LEVEL_MOVE
	COMMAND_LEVEL_STEP
	COMMAND_LEVEL_STOP
	COMMAND_LEVEL_MOVE_TO_LEVEL_WITH_ON_OFF
	COMMAND_LEVEL_MOVE_WITH_ON_OFF
	COMMAND_LEVEL_STEP_WITH_ON_OFF
	COMMAND_LEVEL_STOP_WITH_ON_OFF

	/*** Level Control Move (Mode) Command values ***/
	LEVEL_MOVE_UP   = 0x00
	LEVEL_MOVE_DOWN = 0x01

	/*** Level Control Step (Mode) Command values ***/
	LEVEL_STEP_UP   = 0x00
	LEVEL_STEP_DOWN = 0x01
)

// 16 bit版本号转换为字符串
// 最高位为 版本是否是测试版本, 余下15位分成5,5,5分别为主版本号,次版本号,修正版本号
// 字符串形式 v1.0.0 Beta
func Version(ver uint16) string {
	isBeta := (ver >> 15) == 0x01
	major := (ver >> 10) & 0x1f
	minor := (ver >> 4) & 0x1f
	fixed := ver & 0x1f

	s := fmt.Sprintf("v%d.%d.%d", major, minor, fixed)
	if isBeta {
		return s + " Beta"
	}
	return s
}

// 十进制编译时间转换为字符串形式
// 2019032810 -> 2019-03-28 10h
func BuildDate(val uint32) string {
	hour := val % 100
	val /= 100
	day := val % 100
	val /= 100
	month := val % 100
	year := val / 100

	return fmt.Sprintf("%04d-%02d-%02d %02dh", year, month, day, hour)
}

func PowerSource(p uint8) string {
	powerSource := []string{"Unknown", "single phase", "three phase", "Battery", "DC source", "Emergency mains"}

	msk := p & 0x7f
	if int(msk) < len(powerSource) {
		return powerSource[msk]
	}
	return powerSource[0]
}
