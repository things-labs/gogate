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

// 配置路径
const (
	SmartAppCfgPath = "conf/smartapp.conf" // 应用配置路径
	UsartCfgPath    = "conf/usart.conf"    // 串口配置
)

// 配置文件
var (
	APPCfg  *ini.File
	UartCfg *ini.File
	watcher *fsnotify.Watcher
)

// CfgInit 配置初始化
func CfgInit() error {
	var err error

	if dir := path.Dir(SmartAppCfgPath); !com.IsExist(dir) {
		os.MkdirAll(dir, os.ModePerm)
		FactorySmartAppCfg()
		FactoryUsartCfg()
	}

	if !com.IsExist(SmartAppCfgPath) {
		FactorySmartAppCfg()
	}

	if APPCfg, err = ini.LooseLoad(SmartAppCfgPath); err != nil {
		if APPCfg, err = ini.Load([]byte(SmartAppDefaultCfg)); err != nil {
			logs.Critical("config load failed,", err)
		}
	}

	if !com.IsExist(UsartCfgPath) {
		FactoryUsartCfg()
	}

	if UartCfg, err = ini.LooseLoad(UsartCfgPath); err != nil {
		if UartCfg, err = ini.Load([]byte(UsartDefaultCfg)); err != nil {
			logs.Critical("config load failed,", err)
		}
	}

	go watch()
	return err
}

// FactorySmartAppCfg 恢复smart app为默认
func FactorySmartAppCfg() error {
	f, err := os.Create(SmartAppCfgPath)
	if err != nil {
		return err
	}
	f.WriteString(SmartAppDefaultCfg)
	f.Close()
	return nil
}

// FactoryUsartCfg 串口恢复到出厂设置
func FactoryUsartCfg() error {
	f, err := os.Create(UsartCfgPath)
	if err != nil {
		return err
	}
	f.WriteString(UsartDefaultCfg)
	f.Close()
	return nil
}

// LogsInit log初始化
func LogsInit() {
	logs.Reset() // 复位日志输出流
	sec, err := APPCfg.GetSection("logs")
	if err != nil { // 使用默认控制台配置
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

	watcher.Add(SmartAppCfgPath)
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
