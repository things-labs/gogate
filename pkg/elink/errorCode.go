package elink

// 通用错误码
const (
	CodeSuccess                            = iota // 成功
	CodeErrSysException                           // 内部异常
	CodeErrSysNotSupport                          // 不支持
	CodeErrSysOperationFailed                     // 操作失败
	CodeErrSysInvalidParameter                    // 无效参数
	CodeErrSysInProcess                           // 正在处理中
	CodeErrSysResourceNotSupport                  // 资源不支持
	CodeErrSysMethodNotSupport                    // 方法不支持
	CodeErrSysResourceMethodNotImplemented        // 资源下的方法未实现
)

// CodeErrorMessageInfo 错误码信息
type CodeErrorMessageInfo struct {
	Detail  string
	Message string
}

var codeErrorMessageList = map[int]CodeErrorMessageInfo{
	CodeSuccess:                  {"iot.Success", "success"},
	CodeErrSysException:          {"iot.system.Exception", "Internal system exception"},
	CodeErrSysNotSupport:         {"iot.system.NotSupport", "Not support"},
	CodeErrSysOperationFailed:    {"iot.system.OperationFailed", "Operation failed"},
	CodeErrSysInvalidParameter:   {"iot.system.InvalidParameter", "Invalid parameter"},
	CodeErrSysInProcess:          {"iot.system.InProcess", "In process"},
	CodeErrSysResourceNotSupport: {"iot.system.resource.NotSupport", "Resource not support"},
	CodeErrSysMethodNotSupport:   {"iot.system.method.NotSupport", "Method not support"},
}

// CodeErrorMessage 根据code返回错误码信息
func CodeErrorMessage(code int) CodeErrorMessageInfo {
	errMsg, ok := codeErrorMessageList[code]
	if !ok {
		errMsg = codeErrorMessageList[CodeErrSysNotSupport]
	}
	return errMsg
}

func RegisterCodeErrorMessage(m map[int]CodeErrorMessageInfo) {
	for k, v := range m {
		codeErrorMessageList[k] = v
	}
}
