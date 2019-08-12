package broad

import (
	"fmt"

	"github.com/thinkgos/easyws"
	"github.com/thinkgos/elink"
	"github.com/thinkgos/gomo/lmax"
	"github.com/thinkgos/memlog"

	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

var _ elink.Provider = (*WsProvider)(nil)

type WsProvider struct {
	sess *easyws.Session
}

type DefaultError struct {
	Topic string `json:"topic"`
}

// 默认错误回复,加在topic
func (this *WsProvider) ErrorDefaultResponse(topic string) error {
	o, err := jsoniter.Marshal(DefaultError{topic})
	if err != nil {
		return fmt.Errorf("websocket %v", err)
	}

	return this.sess.WriteMessage(websocket.BinaryMessage, o)
}

// 应答信息
func (this *WsProvider) WriteResponse(tp string, data interface{}) error {
	return this.sess.WriteMessage(websocket.BinaryMessage, data)
}

type wsConsume struct {
	*easyws.Hub
	L *lmax.Lmax
}

func (this *wsConsume) Consume(lower, upper int64) {
	for seq := lower; seq <= upper; seq++ {
		msg := this.L.RingBuffer[seq&lmax.RingBufferDefaultMask]

		err := this.Hub.BroadCast(websocket.BinaryMessage, msg.Data)
		if err != nil {
			memlog.Debug(err)
		}
	}
}
