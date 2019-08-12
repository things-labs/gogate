package misc

import (
	"os"
	"path"

	"github.com/go-ini/ini"
	"github.com/thinkgos/memlog"
	"github.com/thinkgos/utils"
)

// 配置路径
const (
	SmartAppCfgPath = "conf/smartapp.conf" // 应用配置路径
)

// 配置文件
var (
	APPConfig *Config
	appCfg    *ini.File
	//watcher   *fsnotify.Watcher
)

type Logs struct {
	Adapter string `ini:"adapter"`
	Level   int    `ini:"level"`
	IsEFCD  bool   `ini:"isEFCD"`
	IsAsync bool   `ini:"isAsync"`
	LogFCD  int    `ini:"logFCD"`
	Net     string `ini:"net"`
	Addr    string `ini:"addr"`
}

type Usart struct {
	Name     string
	BaudRate int
	DataBit  int
	Parity   string
	StopBit  int
	FlowType int
}

type Config struct {
	OrmDbLog bool  `ini:"ormDbLog"`
	Logs     Logs  `ini:"memlog"`
	Com0     Usart `ini:"com0"`
}

// CfgInit 配置初始化
func CfgInit() {
	var err error

	if dir := path.Dir(SmartAppCfgPath); !utils.IsExist(dir) {
		_ = os.MkdirAll(dir, os.ModePerm)
	}

	if appCfg, err = ini.LooseLoad(SmartAppCfgPath); err != nil {
		appCfg = ini.Empty()
	}
	APPConfig = NewWithDefaultConfig()
	if err := appCfg.MapTo(APPConfig); err != nil {
		memlog.Error(err)
	}
}

// NewWithDefaultConfig 创建一个带默认值的配置
func NewWithDefaultConfig() *Config {
	return &Config{
		OrmDbLog: false,
		Logs: Logs{
			Adapter: memlog.AdapterConsole,
			Level:   7,
			IsEFCD:  false,
			IsAsync: false,
			LogFCD:  0,
			Net:     "udp",
			Addr:    "127.0.0.1:9000",
		},
		Com0: Usart{
			Name:     "/dev/ttyS1",
			BaudRate: 115200,
			DataBit:  8,
			Parity:   "N",
			StopBit:  1,
			FlowType: 0,
		},
	}
}

// SaveConfig 保存本地
func SaveConfig() error {
	cfg := ini.Empty()
	if err := ini.ReflectFrom(cfg, APPConfig); err != nil {
		return err
	}
	return cfg.SaveTo(SmartAppCfgPath)
}
