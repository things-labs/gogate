package broad

import (
	"context"

	"github.com/thinkgos/easyws"
	"github.com/thinkgos/gomo/elink"

	jsoniter "github.com/json-iterator/go"
)

func NewWsHub() *easyws.Hub {
	opt := easyws.NewOptions()
	opt.SetReceiveHandler(func(sess *easyws.Session, t int, data []byte) {
		elink.Server(&WsProvider{sess}, jsoniter.Get(data, "topic").ToString(), data)
	})

	hub := easyws.New(opt)
	go hub.Run(context.TODO())
	return hub
}
