package misc

import (
	"fmt"
	"os"
	"path"

	"github.com/Unknwon/com"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/fsnotify/fsnotify"
	"github.com/go-ini/ini"
)

const (
	SMARTAPP_CFG_PATH = "conf/smartapp.conf" // 应用配置路径
	USART_CFG_PATH    = "conf/usart.conf"    // 串口配置
)

var (
	APPCfg  *ini.File
	UartCfg *ini.File
	watcher *fsnotify.Watcher
)

func CfgInit() error {
	var err error

	if dir := path.Dir(SMARTAPP_CFG_PATH); !com.IsExist(dir) {
		os.MkdirAll(dir, os.ModePerm)
		FactorySmartAppCfg()
		FactoryUsartCfg()
	}

	if !com.IsExist(SMARTAPP_CFG_PATH) {
		FactorySmartAppCfg()
	}

	if APPCfg, err = ini.LooseLoad(SMARTAPP_CFG_PATH); err != nil {
		if APPCfg, err = ini.Load([]byte(SMARTAPP_DEFAULT_CFG)); err != nil {
			logs.Critical("config load failed,", err)
		}
	}

	if !com.IsExist(USART_CFG_PATH) {
		FactoryUsartCfg()
	}

	if UartCfg, err = ini.LooseLoad(USART_CFG_PATH); err != nil {
		if UartCfg, err = ini.Load([]byte(USART_DEFAULT_CFG)); err != nil {
			logs.Critical("config load failed,", err)
		}
	}

	go watch()
	return err
}

func FactorySmartAppCfg() error {
	f, err := os.Create(SMARTAPP_CFG_PATH)
	if err != nil {
		return err
	}
	f.WriteString(SMARTAPP_DEFAULT_CFG)
	f.Close()
	return nil
}

func FactoryUsartCfg() error {
	f, err := os.Create(USART_CFG_PATH)
	if err != nil {
		return err
	}
	f.WriteString(USART_DEFAULT_CFG)
	f.Close()
	return nil
}

func LogsInit() {
	logs.Reset() // 复位日志输出流
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
	// Enalbe async output log
	if sec.Key("isAsync").MustBool(false) {
		logs.Async()
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

	watcher.Add(SMARTAPP_CFG_PATH)
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
			logs.Debug("op:", event)
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
