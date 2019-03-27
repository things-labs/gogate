package mq

import (
	"fmt"
	"time"

	"github.com/slzm40/gomo/elink"
	"github.com/slzm40/gomo/misc"
	"github.com/slzm40/gomo/protocol/elinkch/ctrl"
	"github.com/slzm40/gomo/protocol/elmodels"

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

var gatewayHeartBeatTopic = fmt.Sprintf("data/0/%s/gateway.heartbeat/patch/time", misc.Mac())
var Client mqtt.Client

func init() {
	elink.RegisterTopicInfo(misc.Mac(), gatewayProductKey) // 注册网关产品Key

	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqtt_broker_address).SetPassword(mqtt_broker_password) // broker
	opts.SetUsername("1").SetClientID(misc.Mac())

	opts.SetCleanSession(false)
	opts.SetAutoReconnect(true)

	opts.SetOnConnectHandler(func(cli mqtt.Client) {
		logs.Info("mqtt client connect success")
		chList := elink.ChannelSelectorList()
		for _, ch := range chList {
			s := fmt.Sprintf("%s/%s/%s/+/+/+/#", ch, elink.TpInfos.ProductKey, misc.Mac())
			cli.Subscribe(s, 2, elink.Server)
		}

		time.AfterFunc(time.Second, HeartBeatStatus)
	})

	opts.SetConnectionLostHandler(func(cli mqtt.Client, err error) {
		logs.Warn("mqtt client connection lost ", err)
	})

	if out, err := jsoniter.Marshal(elmodels.GatewayHeatbeats(false)); err != nil {
		logs.Error("mqtt %s", err.Error())
	} else {
		opts.SetBinaryWill(gatewayHeartBeatTopic, out, 2, false)
	}
	Client = mqtt.NewClient(opts)
	started()
}

func started() {
	logs.Info("mqtt client connecting...")
	if token := Client.Connect(); token.Wait() && token.Error() != nil {
		logs.Error("mqtt client connect failed, ", token.Error())
		time.AfterFunc(time.Second*30, started)
	}
}

func HeartBeatStatus() {
	defer time.AfterFunc(time.Second*30, HeartBeatStatus)
	if !Client.IsConnected() {
		return
	}

	out, err := jsoniter.Marshal(elmodels.GatewayHeatbeats(true))
	if err != nil {
		logs.Error("HeartBeatStatus:", err)
		return
	}
	Client.Publish(gatewayHeartBeatTopic, 0, false, out)

}

func WritePublishChData(resourse, method, messageType string, data interface{}) error {
	return elink.WritePublish(Client, ctrl.ChannelData, resourse, method, messageType, data)
}
