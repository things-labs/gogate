package broad

import (
	"time"

	"github.com/thinkgos/easyws"

	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
	"github.com/thinkgos/gogate/misc"
	"github.com/thinkgos/gogate/protocol/elinkmd"
	"github.com/thinkgos/gomo/elink"
	"github.com/thinkgos/gomo/lmax"
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

// 网关心跳包
func HeartBeatStatus() {
	defer time.AfterFunc(HeartBeatTime, HeartBeatStatus)

	// 心跳包推送
	func() {
		tp := elink.FormatPshTopic(elink.ChannelInternal,
			elinkmd.GatewayHeartbeat, elink.MethodPatch, elink.MessageTypeTime)
		out, err := jsoniter.Marshal(elinkmd.GatewayHeatbeats(tp, true))
		if err != nil {
			logs.Error("GatewayHeatbeats:", err)
			return
		}
		err = Disrup.Publish(tp, out)
		if err != nil {
			logs.Error("GatewayHeatbeats:", err)
		}
	}()

	// 系统监控信息推送
	func() {
		tp := elink.FormatPshTopic(elink.ChannelInternal,
			elinkmd.SystemMonitor, elink.MethodPatch, elink.MessageTypeTime)

		out, err := jsoniter.Marshal(elinkmd.GatewayMonitors(tp))
		if err != nil {
			logs.Error("GatewayMonitors:", err)
			return
		}
		err = Disrup.Publish(tp, out)
		if err != nil {
			logs.Error("GatewayHeatbeats:", err)
		}
	}()
}

func Publish(tp string, data interface{}) error {
	return Disrup.Publish(tp, data)
}
