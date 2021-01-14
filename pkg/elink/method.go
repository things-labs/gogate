package elink

import (
	"strings"
)

// 请求的方法
const (
	MethodUnknown = "unknown"
	MethodGet     = "get"
	MethodPost    = "post"
	MethodPut     = "put"
	MethodDelete  = "delete"
)

// 消息类型
const (
	MessageTypeAnnce = "annce"
	MessageTypeTime  = "time"
	MessageTypeAlarm = "alarm"
)

var (
	// 支持的方法
	MethodSupport = []string{MethodGet, MethodPost, MethodPut, MethodDelete}
	// 支持的消息类型
	MessageTypeSupport = []string{MessageTypeAnnce, MessageTypeTime, MessageTypeAlarm}
)

// HasMethod 是否有相应的请求方法, 是true,否则false
func HasMethod(method string) bool {
	for _, s := range MethodSupport {
		if strings.EqualFold(s, method) {
			return true
		}
	}
	return false
}

func ContainsUnknownMethod(str []string) bool {
	for _, s := range str {
		if strings.EqualFold(s, MethodUnknown) {
			return true
		}
	}
	return false
}

// HasMessageType 判断有elink消息类型, 是true,否则false
func HasMessageType(msgType string) bool {
	for _, s := range MessageTypeSupport {
		if strings.EqualFold(s, msgType) {
			return true
		}
	}
	return false
}
