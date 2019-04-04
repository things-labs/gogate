package mq

import (
	"fmt"
	"sync"
	"time"

	"github.com/thinkgos/gogate/protocol/elinkch/ctrl"
	"github.com/thinkgos/gogate/protocol/elinkmd"
	"github.com/thinkgos/gogate/protocol/elinkres"
	"github.com/thinkgos/gomo/elink"
	"github.com/thinkgos/gomo/misc"

	"github.com/astaxie/beego/logs"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/json-iterator/go"
)

const (
	mqtt_broker_address  = "mqtt.lchtime.cn:1883"
	mqtt_broker_password = "52399399"
)

const (
	gatewayProductKey = "lc_gzs100"
)

var Client mqtt.Client
var heartOnce sync.Once

func init() {
	elink.RegisterTopicInfo(misc.Mac(), gatewayProductKey) // 注册网关产品Key

	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqtt_broker_address).SetClientID(misc.Mac()) // broker and clientID
	opts.SetUsername("1").SetPassword(mqtt_broker_password)     // user name and password
	opts.SetCleanSession(false).SetAutoReconnect(true)

	opts.SetOnConnectHandler(func(cli mqtt.Client) {
		logs.Info("mqtt client connect success")
		chList := elink.ChannelSelectorList()
		for _, ch := range chList {
			s := fmt.Sprintf("%s/%s/%s/+/+/+/#", ch, elink.TpInfos.ProductKey, misc.Mac())
			cli.Subscribe(s, 2, elink.Server)
		}
		heartOnce.Do(func() { time.AfterFunc(time.Second, HeartBeatStatus) })
	})

	opts.SetConnectionLostHandler(func(cli mqtt.Client, err error) {
		logs.Warn("mqtt client connection lost, ", err)
	})

	if out, err := jsoniter.Marshal(elinkmd.GatewayHeatbeats(false)); err != nil {
		logs.Error("mqtt %s", err.Error())
	} else {
		opts.SetBinaryWill(
			fmt.Sprintf("data/0/%s/%s/patch/time", misc.Mac(), elinkres.GatewayHeartbeat),
			out, 2, false)
	}
	Client = mqtt.NewClient(opts)
	started()
}

// 启动连接mqtt
func started() {
	logs.Info("mqtt client connecting...")
	if token := Client.Connect(); token.Wait() && token.Error() != nil {
		logs.Error("mqtt client connect failed, ", token.Error())
		time.AfterFunc(time.Second*30, started)
	}
}

// 网关心跳包
func HeartBeatStatus() {
	defer time.AfterFunc(time.Second*30, HeartBeatStatus)
	if !Client.IsConnected() {
		return
	}

	func() {
		out, err := jsoniter.Marshal(elinkmd.GatewayHeatbeats(true))
		if err != nil {
			logs.Error("GatewayHeatbeats:", err)
			return
		}
		elink.WriteSpecialData(Client, ctrl.ChannelData,
			elinkres.GatewayHeartbeat, elink.MethodPatch, elink.MessageTypeTime, out)
	}()

	func() {
		out, err := jsoniter.Marshal(elinkmd.GatewayMonitors())
		if err != nil {
			logs.Error("GatewayMonitors:", err)
			return
		}
		elink.WriteSpecialData(Client, ctrl.ChannelData,
			elinkres.GatewayMonitor, elink.MethodPatch, elink.MessageTypeTime, out)
	}()

}

// ctrl data通道推送数据
func WriteCtrlData(resourse, method, messageType string, payload []byte) error {
	return ctrl.WriteData(Client, resourse, method, messageType, payload)
}
