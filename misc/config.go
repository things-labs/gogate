package misc

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/fsnotify/fsnotify"
	"github.com/go-ini/ini"
)

const (
	APP_CFG_PATH   = "./conf/moapp.conf" // 应用配置路径
	USART_CFG_PATH = "./conf/usart.conf" // 串口配置
)

var (
	APPCfg  *ini.File
	UartCfg *ini.File
	watcher *fsnotify.Watcher
)

func init() {
	var err error

	if APPCfg, err = ini.Load(APP_CFG_PATH); err != nil {
		panic(err)
	}

	if UartCfg, err = ini.Load(USART_CFG_PATH); err != nil {
		panic(err)
	}

	go watch()
}

func LogsInit() {
	logs.Reset() // 复位日志
	sec, err := APPCfg.GetSection("logs")
	if err != nil {
		level := 7
		if beego.BConfig.RunMode == "prod" {
			level = 3
		}
		logs.SetLogger(logs.AdapterConsole,
			fmt.Sprintf(`{"level":%d,"color":true}`, level)) // out to console
		logs.Debug("use adapter: console,level: %d", level)
		return
	}
	adapter := sec.Key("adapter").MustString("console")
	tmpll := sec.Key("level").MustInt(7)
	if adapter == logs.AdapterConn {
		/* default: {"net":"udp","addr":"127.0.0.1:8080",
		"level":7,"reconnect":true,"color":true} */
		tmpnet := sec.Key("net").MustString("udp")
		tmpaddr := sec.Key("addr").MustString("127.0.0.1:8080")
		logs.SetLogger(logs.AdapterConn,
			fmt.Sprintf(`{"net":"%s","addr":"%s","level":%d,"reconnect":true,"color":true}`,
				tmpnet, tmpaddr, tmpll))
	} else {
		logs.SetLogger(logs.AdapterConsole,
			fmt.Sprintf(`{"level":%d,"color":true}`, tmpll)) // out to console
	}
	// Enable output filename and line
	logs.EnableFuncCallDepth(sec.Key("isEFCD").MustBool(false))
	if sec.Key("isAsync").MustBool(false) {
		logs.Async() // Enalbe async output log
	}
	logs.Debug("use adapter: %s,level: %d", adapter, tmpll)
	return
}

func watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()

	watcher.Add(APP_CFG_PATH)
	for {
		select {
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logs.Debug("fsnotify:", err)
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if (event.Op & fsnotify.Write) == fsnotify.Write {
				err = APPCfg.Reload()
				if err != nil {
					logs.Debug("config reload failed!", err)
					break
				}
				logs.Info("wirte happen")
				LogsInit()
			}
		}
	}
}
