package elinkmd

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/thinkgos/gomo/elink"
)

type ItemInfos struct {
	Client    mqtt.Client
	Tp        *elink.TopicLayer
	ProductID int
	Sn        string
	Pkid      int
	IsLocal   bool // if local do not send message to up
}
