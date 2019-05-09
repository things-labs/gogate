package broad

import (
	"time"

	"github.com/thinkgos/easyws"
	"github.com/thinkgos/gogate/apps/elinkch/ctrl"
	"github.com/thinkgos/gogate/apps/elinkmd"
	"github.com/thinkgos/gogate/misc"
	"github.com/thinkgos/gomo/elink"
	"github.com/thinkgos/gomo/lmax"

	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

const (
	HeartBeatTime = 30 * time.Second
)

var WsHub *easyws.Hub
var Disrup *lmax.Lmax

func BroadInit() {
	Disrup = lmax.New()
	mqCli := NewMqClient(elinkmd.ProductKey, misc.Mac())
	WsHub = NewWsHub()
	go Disrup.Run(&mqConsume{mqCli, Disrup}, &wsConsume{WsHub, Disrup})
	time.Sleep(time.Millisecond * 1)
	HeartBeatStatus()
}

func PublishServerJSON(tp string, data interface{}) error {
	out, err := jsoniter.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "json marshal failed")
	}
	return Disrup.Publish(tp, out)
}

// PublishPyServerJSON 推送数据,通道推送数据
func PublishPyServerJSON(tp string, payload interface{}) error {
	return PublishServerJSON(tp, &ctrl.PublishData{&ctrl.BasePublishData{tp}, payload})
}

// HeartBeatStatus 网关心跳包
func HeartBeatStatus() {
	defer time.AfterFunc(HeartBeatTime, HeartBeatStatus)

	// 心跳包推送
	func() {
		tp := elink.FormatPshTopic(elink.ChannelInternal, elinkmd.GatewayHeartbeat,
			elink.MethodPatch, elink.MessageTypeTime)
		err := PublishPyServerJSON(tp, elinkmd.GatewayHeatbeats(true))
		if err != nil {
			logs.Error("GatewayHeatbeats:", err)
		}
	}()

	// 系统监控信息推送
	func() {
		gm, err := elinkmd.GatewayMonitors()
		if err != nil {
			logs.Error("GatewayMonitors:", err)
			return
		}
		tp := elink.FormatPshTopic(elink.ChannelInternal, elinkmd.SystemMonitor,
			elink.MethodPatch, elink.MessageTypeTime)
		err = PublishPyServerJSON(tp, gm)
		if err != nil {
			logs.Error("GatewayMonitors:", err)
		}
	}()
}
