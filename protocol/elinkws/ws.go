package elinkws

import (
	"context"
	"net"
	"time"

	"github.com/thinkgos/gogate/protocol/elinkch/ctrl"
	"github.com/thinkgos/gomo/elink"

	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

const (
	sendSize     = 32
	writeWait    = 1 * time.Second
	keepAlive    = 60 * time.Second
	monitorAlive = keepAlive * 110 / 100
)

var _ elink.Provider = (*Provider)(nil)

type Provider struct {
	Conn  *websocket.Conn
	send  chan []byte
	alive chan struct{}
}

// 创建mqtt provider实例
func NewProvider(c *websocket.Conn) *Provider {
	// ctx, cancel := context.WithCancel(context.Background())
	return &Provider{c,
		make(chan []byte, sendSize),
		make(chan struct{}, 1),
	}
}

// 默认错误回误,加在topic
func (this *Provider) ErrorDefaultResponse(topic string) error {
	o, err := jsoniter.Marshal(ctrl.BaseData{topic})
	if err != nil {
		return errors.Wrap(err, "websocket")
	}
	this.send <- o
	return nil
}

// 应答信息
func (this *Provider) WriteResponse(tp string, data interface{}) error {
	return this.Publish(tp, data)
}

// 数据推送
func (this *Provider) Publish(tp string, data interface{}) error {
	var py []byte

	switch data.(type) {
	case string:
		py = []byte(data.(string))
	case []byte:
		py = data.([]byte)
	default:
		return errors.New("Unknown data type")
	}
	this.send <- py
	return nil
}

func (this *Provider) Run(ctx context.Context) {
	client := elink.NewClient(elink.Hub, this)
	client.Hub.ManagaClient(true, client)

	this.Conn.SetPongHandler(func(string) error {
		logs.Debug("%s pong", this.Conn.RemoteAddr().String())
		this.alive <- struct{}{}
		return nil
	})
	this.Conn.SetPingHandler(func(message string) error {
		err := this.Conn.WriteControl(websocket.PongMessage,
			[]byte(message), time.Now().Add(writeWait))
		if err != nil {
			if err == websocket.ErrCloseSent {
				// see default handler
			} else if e, ok := err.(net.Error); ok && e.Temporary() {
				// see default handler
			} else {
				return err
			}
		}
		logs.Debug("%s ping", this.Conn.RemoteAddr().String())
		this.alive <- struct{}{}
		return nil
	})

	lctx, cancel := context.WithCancel(ctx)
	closeFunc := func() error {
		client.Hub.ManagaClient(false, client)
		return this.Conn.Close()
	}

	go func() {
		var retries int

		monTick := time.NewTimer(monitorAlive)
		defer func() {
			logs.Error("Run write: closed")
			closeFunc()
		}()
		for {
			select {
			case <-lctx.Done():
				return
			case msg, ok := <-this.send:
				this.Conn.SetWriteDeadline(time.Now().Add(writeWait))
				if !ok {
					this.Conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				err := this.Conn.WriteMessage(websocket.BinaryMessage, msg)
				if err != nil {
					logs.Error("Run write: ", err)
					return
				}

			case <-this.alive:
				retries = 0
				monTick.Reset(monitorAlive)

			case <-monTick.C:
				if retries++; retries > 3 {
					monTick.Stop()
					return
				}
				monTick.Reset(monitorAlive / 2)
				err := this.Conn.WriteControl(websocket.PingMessage, []byte{},
					time.Now().Add(writeWait))
				if err != nil {
					logs.Error("server Write: ", err)
					return
				}
			}
		}
	}()

	defer func() {
		cancel()
		closeFunc()
	}()

	for {
		_, msg, err := this.Conn.ReadMessage()
		if err != nil {
			logs.Error("Run Read: ", err)
			break
		}
		elink.Server(this, jsoniter.Get(msg, "topic").ToString(), msg)
	}
}
