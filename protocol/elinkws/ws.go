package elinkws

import (
	"context"
	"net"
	"sync/atomic"
	"time"

	"github.com/thinkgos/gogate/protocol/elinkch/ctrl"
	"github.com/thinkgos/gomo/elink"

	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

var _ elink.Provider = (*Provider)(nil)

type Provider struct {
	Conn     *websocket.Conn
	cfg      *Config
	outBound chan []byte
	alive    int32
}

// 创建mqtt provider实例
func NewProvider(c *websocket.Conn, cfg ...*Config) *Provider {
	config := newDefaultConfig()
	if len(cfg) > 0 {
		if cfg[0].Radtio < DefaultRadtio {
			cfg[0].Radtio = DefaultRadtio
		}
		config = cfg[0]
	}
	return &Provider{
		c,
		config,
		make(chan []byte, config.MaxMessageSize),
		0,
	}
}

// 默认错误回误,加在topic
func (this *Provider) ErrorDefaultResponse(topic string) error {
	o, err := jsoniter.Marshal(ctrl.BaseData{topic})
	if err != nil {
		return errors.Wrap(err, "websocket")
	}
	this.outBound <- o
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
	this.outBound <- py
	return nil
}

func (this *Provider) writeDump(ctx context.Context) {
	var retries int

	cfg := this.cfg
	monTick := time.NewTicker(cfg.KeepAlive * time.Duration(cfg.Radtio) / 100)
	defer func() {
		logs.Error("Run write: closed")
		monTick.Stop()
		this.Conn.Close()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-this.outBound:
			this.Conn.SetWriteDeadline(time.Now().Add(cfg.WriteWait))
			if !ok {
				this.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := this.Conn.WriteMessage(websocket.BinaryMessage, msg)
			if err != nil {
				logs.Error("Run write: ", err)
				return
			}
		case <-monTick.C:
			if atomic.AddInt32(&this.alive, 1) > 1 {
				if retries++; retries > 3 {
					return
				}
				err := this.Conn.WriteControl(websocket.PingMessage, []byte{},
					time.Now().Add(cfg.WriteWait))
				if err != nil {
					logs.Error("server Write: ", err)
					return
				}
			} else {
				retries = 0
			}
		}
	}
}

func (this *Provider) Run(ctx context.Context) {
	client := elink.NewClient(elink.Hub, this)
	client.Hub.ManagaClient(true, client)

	lctx, cancel := context.WithCancel(ctx)
	go this.writeDump(lctx)

	readWait := this.cfg.KeepAlive * time.Duration(this.cfg.Radtio) /
		100 * (tuple + 1)

	this.Conn.SetPongHandler(func(string) error {

		atomic.StoreInt32(&this.alive, 0)
		this.Conn.SetReadDeadline(time.Now().Add(readWait))
		logs.Debug("%s pong", this.Conn.RemoteAddr().String())
		return nil
	})
	this.Conn.SetPingHandler(func(message string) error {
		atomic.StoreInt32(&this.alive, 0)
		this.Conn.SetReadDeadline(time.Now().Add(readWait))
		err := this.Conn.WriteControl(websocket.PongMessage,
			[]byte(message), time.Now().Add(this.cfg.WriteWait))
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
		return nil
	})

	if this.cfg.MaxMessageSize > 0 {
		this.Conn.SetReadLimit(this.cfg.MaxMessageSize)
	}
	this.Conn.SetReadDeadline(time.Now().Add(readWait))
	for {
		_, msg, err := this.Conn.ReadMessage()
		if err != nil {
			logs.Error("Run Read: ", err)
			break
		}
		elink.Server(this, jsoniter.Get(msg, "topic").ToString(), msg)
	}

	client.Hub.ManagaClient(false, client)
	this.Conn.Close()
	cancel()
}
