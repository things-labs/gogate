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
)

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
