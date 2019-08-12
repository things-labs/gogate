package broad

import (
	"github.com/thinkgos/elink"
	"github.com/thinkgos/gomo/lmax"
	"github.com/thinkgos/memlog"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var _ elink.Provider = (*MqProvider)(nil)

type MqProvider struct {
	cli mqtt.Client
}

// 错误加在主题上的回复
func (this *MqProvider) ErrorDefaultResponse(topic string) error {
	return this.WriteResponse(topic, "{}")
}

// 应答信息
func (this *MqProvider) WriteResponse(topic string, data interface{}) error {
	return this.cli.Publish(topic, 2, false, data).Error()
}

// 回调
func MessageHandle(client mqtt.Client, message mqtt.Message) {
	memlog.Debug("Topic: %s", message.Topic())
	memlog.Warn("MessageID: %d,Qos - %d,Retained - %t,Duplicate - %t",
		message.MessageID(), message.Qos(), message.Retained(), message.Duplicate())
	//memlog.Debug("receive:\n%s\n", message.Payload())

	// 抛弃retain 和重复的消息 必须使用Qos = 2的消息
	if message.Retained() || message.Duplicate() || message.Qos() != 2 {
		memlog.Warn("Handle: Invalid message discard")
		return
	}
	elink.Server(&MqProvider{client}, message.Topic(), message.Payload())
}

type mqConsume struct {
	mqtt.Client
	L *lmax.Lmax
}

func (this *mqConsume) Consume(lower, upper int64) {
	for seq := lower; seq <= upper; seq++ {
		msg := this.L.RingBuffer[seq&lmax.RingBufferDefaultMask]
		if err := this.Client.Publish(msg.Topic, 1, false, msg.Data).Error(); err != nil {
			memlog.Debug(err)
		}
	}
}
