package elinkws

import (
	"context"
	"errors"
	"time"

	"github.com/thinkgos/gogate/protocol/elinkch/ctrl"

	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/thinkgos/gomo/elink"
)

var _ elink.Provider = (*Provider)(nil)

type Provider struct {
	C      *websocket.Conn
	ctx    context.Context
	cancel context.CancelFunc

	send chan []byte
}

type message struct {
	*ctrl.BaseResponse
	*ctrl.BaseRawPayload
}

// 创建mqtt provider实例
func NewProvider(c interface{}) elink.Provider {
	ctx, cancel := context.WithCancel(context.Background())
	return &Provider{c.(*websocket.Conn),
		ctx,
		cancel,
		make(chan []byte, 32),
	}
}

// 应答信息
func (this *Provider) WriteResponse(tp string, data interface{}) error {
	var py []byte

	switch data.(type) {
	case string:
		py = []byte(data.(string))
	case []byte:
		py = data.([]byte)
	default:
		return errors.New("Unknown payload type")
	}

	rsp := &struct {
		*ctrl.BaseResponse
		*ctrl.BaseRawPayload
	}{}

	err := jsoniter.Unmarshal(py, rsp)
	if err != nil {
		return err
	}
	rsp.Topic = tp
	o, err := jsoniter.Marshal(rsp)
	if err != nil {
		return err
	}

	this.send <- o
	return nil
}

// 数据推送
func (this *Provider) Publish(tp string, data interface{}) error {
	this.send <- []byte(tp)
	return nil
}

func (this *Provider) RunWrite(client *elink.Client) {
	defer func() {
		this.C.Close()
		elink.ManagaClient(false, client)
	}()
	for {
		select {
		case msg, ok := <-this.send:
			this.C.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				this.C.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := this.C.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		}
	}
}

func (this *Provider) RunRead(client *elink.Client) {
	for {
		_, msg, err := this.C.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logs.Warn("RunRead: %v", err)
			}
			break
		}
		tp := jsoniter.Get(msg, "topic").ToString()
		if len(tp) == 0 {
			logs.Warn("Handle: Invalid topic discard")
			continue
		}
		elink.Server(NewProvider(this.C), tp, msg)
	}

	this.C.Close()
	elink.ManagaClient(false, client)
}
