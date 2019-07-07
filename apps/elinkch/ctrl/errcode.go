package ctrl

import (
	"github.com/thinkgos/elink"
)

const (
	CodeErrCommonUserNoAccess                             = 50 + iota // 用户无权限
	CodeErrCommonAuthorizationSignatureVerificationFailed             // 签名验证失败
)

// 用户定义错误码
const (
	CodeErrProudctUndefined             = 100 + iota // 产品未定义
	CodeErrProudctOperationFailed                    // 产品
	CodeErrProudctFeatureUndefined                   // 产品功能未定义
	CodeErrDeviceOffline                             // 设备离线
	CodeErrDeviceNotExist                            // 无此设备
	CodeErrDeviceOperationFailed                     // 设备操作失败
	CodeErrDeviceFeatureNotSupport                   // 设备功能不支持
	CodeErrDeviceCommandNotSupport                   // 设备命令不支持
	CodeErrDeviceCommandOperationFailed              // 设备命令操作失败
	CodeErrDevicePropertysNotSupport                 // 设备属性不支持
)

func init() {
	elink.RegisterCodeErrorMessage(map[int]elink.CodeErrorMessageInfo{
		// 通用错误码
		CodeErrCommonUserNoAccess:                             {"iot.common.user.NoAccess", "No access"},
		CodeErrCommonAuthorizationSignatureVerificationFailed: {"iot.common.authorization.SignatureverificationFailed", "Signature verification failed"},
		// 产品错误码
		CodeErrProudctUndefined:        {"iot.product.Undefined", "Product undefined"},
		CodeErrProudctOperationFailed:  {"iot.product.OperationFailed", "Product operation failure"},
		CodeErrProudctFeatureUndefined: {"iot.product.FeatureUndefined", "Product feature not support"},
		// 设备错误码
		CodeErrDeviceOffline:                {"iot.device.Offline", "Device offline"},
		CodeErrDeviceNotExist:               {"iot.device.NotExist", "Device not exist"},
		CodeErrDeviceOperationFailed:        {"iot.device.OperationFailed", "Device operation failure"},
		CodeErrDeviceFeatureNotSupport:      {"iot.device.FeatureNotSupport", "Device feature not support"},
		CodeErrDeviceCommandNotSupport:      {"iot.device.CommandNotSupport", "Device command not support"},
		CodeErrDeviceCommandOperationFailed: {"iot.device.CommandOperationFailed", "Device command operation failure"},
		CodeErrDevicePropertysNotSupport:    {"iot.device.PropertysNotSupport", "Device propertys not support"},
		// 用户定义级
	})
}
