package elmodels

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/thinkgos/gomo/elink"
)

type BaseSnPayload struct {
	ProductID int    `json:"productID"`
	Sn        string `json:"sn"`
}

type DevicesInfo struct {
	ProductID int      `json:"productID"`
	Sn        []string `json:"sn"`
}

type ItemInfos struct {
	Client    mqtt.Client
	Tp        *elink.TopicLayer
	ProductID int
	Sn        string
	Pkid      int
	IsLocal   bool // if local do not send message to up
}
