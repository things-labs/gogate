package main

import (
	"fmt"
	"time"

	_ "github.com/slzm40/gomo/elink/channel/ctrl"

	"github.com/astaxie/beego/logs"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/json-iterator/go"
	"github.com/slzm40/common"
	"github.com/slzm40/gomo/elink"
	"github.com/slzm40/gomo/misc"
)

const (
	mqtt_broker_address  = "mqtt.lchtime.cn:1883"
	mqtt_broker_password = "52399399"
)

const (
	gatewayProductKey = "zbgw01"
)

var MqClinet mqtt.Client

func init() {
	elink.RegisterTopicInfo(misc.Mac(), gatewayProductKey) // 注册网关产品Key

	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqtt_broker_address) // broker
	opts.SetPassword(mqtt_broker_password)
	opts.SetUsername("1")
	opts.SetAutoReconnect(true)

	opts.SetClientID(misc.Mac())
	opts.SetCleanSession(false)
	opts.SetOnConnectHandler(func(cli mqtt.Client) {
		logs.Info("mqtt client connect success")
		chList := elink.ChannelSelectorList()
		for _, ch := range chList {
			s := fmt.Sprintf("%s/%s/%s/+/+/+/#", ch, elink.TpInfos.ProductKey, misc.Mac())
			cli.Subscribe(s, 2, elink.Server)
		}

		time.AfterFunc(time.Second*10, HeartBeatStatus)
	})

	opts.SetConnectionLostHandler(func(cli mqtt.Client, err error) {
		logs.Warn("mqtt clinet connection lost ", err)
	})
	MqClinet = mqtt.NewClient(opts)
	started()
}

func started() {
	logs.Info("mqtt client connecting...")
	if token := MqClinet.Connect(); token.Wait() && token.Error() != nil {
		logs.Error("mqtt client connect failed, ", token.Error())
		time.AfterFunc(time.Second*30, started)
	}
}

type DeviceInfo struct {
	Sn string `json:"sn"`
}
type DeviceStatus struct {
	CurrentTime   string `json:"currentTime"`
	StartDateTime string `json:"startDateTime"`
	RunningTime   string `json:"runningTime"`
	Status        string `json:"status"`
}
type NetInfo struct {
	MAC string `json:"MAC"`
	Mac string `json:"mac"`
}

type Info struct {
	DeviceInfo   DeviceInfo   `json:"device_info"`
	DeviceStatus DeviceStatus `json:"device_status"`
	NetInfo      NetInfo      `json:"net_info"`
}

type GatewayStatus struct {
	Info Info `json:"info"`
}

func HeartBeatStatus() {
	defer time.AfterFunc(time.Second*10, HeartBeatStatus)
	if !MqClinet.IsConnected() {
		return
	}

	mac := misc.Mac()

	dInfo := DeviceInfo{Sn: mac}
	dStatus := DeviceStatus{
		CurrentTime:   time.Now().Local().Format("2006-01-02 15:04:05"),
		StartDateTime: common.SetupTime(),
		RunningTime:   common.RunningTime(),
		Status:        "online",
	}
	nInfo := NetInfo{
		MAC: misc.MAC(),
		Mac: mac,
	}

	gStatus := GatewayStatus{
		Info: Info{
			DeviceInfo:   dInfo,
			DeviceStatus: dStatus,
			NetInfo:      nInfo,
		},
	}

	out, err := jsoniter.Marshal(gStatus)
	if err != nil {
		logs.Error("HeartBeatStatus:", err)
	} else {
		s := fmt.Sprintf("data/0/%s/gateway.heartbeat/patch/time", mac)
		MqClinet.Publish(s, 0, false, out)
	}
}
