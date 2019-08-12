package broad

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	jsoniter "github.com/json-iterator/go"
	"github.com/thinkgos/elink"
	"github.com/thinkgos/memlog"

	"github.com/thinkgos/gogate/apps/elinkch/ctrl"
	"github.com/thinkgos/gogate/apps/elinkmd"
	"github.com/thinkgos/gogate/misc"
)

const (
	HeartBeatTime = 30 * time.Second
)

var mqCli mqtt.Client

func BroadInit() {
	mqCli = NewMqClient(elinkmd.ProductKey, misc.Mac())
	time.Sleep(time.Millisecond * 1)
	HeartBeatStatus()
}

// PublishServerJSON 推送json数据
func PublishServerJSON(tp string, data interface{}) error {
	out, err := jsoniter.Marshal(data)
	if err != nil {
		return err
	}
	return mqCli.Publish(tp, 1, false, out).Error()
}

// PublishPyServerJSON 推送数据,通道推送数据,只给出负载
func PublishPyServerJSON(tp string, payload interface{}) error {
	return PublishServerJSON(tp, &ctrl.PublishData{
		BasePublishData: &ctrl.BasePublishData{Topic: tp},
		Payload:         payload})
}

// HeartBeatStatus 网关心跳包
func HeartBeatStatus() {
	defer time.AfterFunc(HeartBeatTime, HeartBeatStatus)

	// 心跳包推送
	func() {
		tp := ctrl.EncodePushTopic(elink.ChannelInternal, elinkmd.GatewayHeartbeat,
			elink.MethodPut, elink.MessageTypeTime)
		err := PublishPyServerJSON(tp, elinkmd.GetGatewayHeatbeatInfo(true))
		if err != nil {
			memlog.Error("GetGatewayHeatbeatInfo:", err)
		}
	}()

	// 系统监控信息推送
	func() {
		gm, err := elinkmd.GetGatewayMonitorInfo()
		if err != nil {
			memlog.Error("GetGatewayMonitorInfo:", err)
			return
		}
		tp := ctrl.EncodePushTopic(elink.ChannelInternal, elinkmd.SystemMonitor,
			elink.MethodPut, elink.MessageTypeTime)
		if err = PublishPyServerJSON(tp, gm); err != nil {
			memlog.Error("GetGatewayMonitorInfo:", err)
		}
	}()
}
