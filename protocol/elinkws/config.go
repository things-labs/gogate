package elinkws

import (
	"time"
)

const (
	tuple              = 3
	DefaultWriteWait   = 1 * time.Second
	DefaultKeepAlive   = 60 * time.Second
	DefaultRadtio      = 110
	DefaultMessageSize = 32
)

// websocket 配置
type Config struct {
	WriteWait         time.Duration // 写超时时间
	KeepAlive         time.Duration // 保活时间
	Radtio            int           // 监控比例, 需大于100,默认系统是110 即比例1.1
	MaxMessageSize    int64         // 消息最大字节数, 如果为0,使用系统默认设置
	MessageBufferSize int           // 消息缓存数
}

func newDefaultConfig() *Config {
	return &Config{
		WriteWait:         DefaultWriteWait,
		KeepAlive:         DefaultKeepAlive,
		Radtio:            DefaultRadtio,
		MaxMessageSize:    0,
		MessageBufferSize: DefaultMessageSize,
	}
}
