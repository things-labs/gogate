package broad

import (
	"context"

	"github.com/thinkgos/gomo/elink"

	jsoniter "github.com/json-iterator/go"
	"github.com/thinkgos/easyws"
)

var WsHub *easyws.Hub

func WsInit() {
	opt := easyws.NewOptions()
	opt.SetReceiveHandler(func(sess *easyws.Session, t int, data []byte) {
		elink.Server(&WsProvider{sess}, jsoniter.Get(data, "topic").ToString(), data)
	})

	WsHub = easyws.New(opt)
	go WsHub.Run(context.TODO())
}
