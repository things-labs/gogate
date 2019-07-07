package misc

import (
	"fmt"
	"os"
	"path"

	"github.com/astaxie/beego/logs"
	"github.com/go-ini/ini"
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
	Logs     Logs  `ini:"logs"`
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
		logs.Error(err)
	}
}

// NewWithDefaultConfig 创建一个带默认值的配置
func NewWithDefaultConfig() *Config {
	return &Config{
		OrmDbLog: false,
		Logs: Logs{
			Adapter: logs.AdapterConsole,
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

// LogsInit log初始化
func LogsInit() {
	logs.Reset() // 复位日志输出流

	if APPConfig.Logs.Adapter == logs.AdapterConn {
		/* default: {"net":"udp","addr":"127.0.0.1:8080","level":7,"reconnect":true,"color":true} */
		_ = logs.SetLogger(logs.AdapterConn,
			fmt.Sprintf(`{"net":"%s","addr":"%s","level":%d,"reconnect":true,"color":true}`,
				APPConfig.Logs.Net, APPConfig.Logs.Addr, APPConfig.Logs.Level))
	} else {
		_ = logs.SetLogger(logs.AdapterConsole,
			fmt.Sprintf(`{"level":%d,"color":true}`, APPConfig.Logs.Level)) // out to console
	}
	// Enable output filename and line
	logs.EnableFuncCallDepth(APPConfig.Logs.IsEFCD)
	// Enalbe async output log
	if APPConfig.Logs.IsAsync {
		logs.Async()
	}

	return
}
