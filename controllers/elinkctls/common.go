package elinkctls

import (
	"github.com/slzm40/gogate/apps"
	"github.com/slzm40/gomo/elink"
	"github.com/slzm40/gomo/elink/channel/ctrl"
)

func WritePublishChData(resourse, method, messageType string, data interface{}) error {
	return elink.WritePublish(apps.MqClinet, ctrl.ChannelData, resourse, method, messageType, data)
}
