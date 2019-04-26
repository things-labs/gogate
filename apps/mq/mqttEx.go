package mq

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"sync"
	"time"

	"github.com/thinkgos/gogate/misc"
	"github.com/thinkgos/gogate/protocol/elinkch/ctrl"
	"github.com/thinkgos/gogate/protocol/elinkmd"
	"github.com/thinkgos/gogate/protocol/elinkmq"
	"github.com/thinkgos/gomo/elink"

	"github.com/astaxie/beego/logs"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	jsoniter "github.com/json-iterator/go"
)

const (
	HeartBeatTime = 10 * time.Second
)

const (
	mqtt_broker_address = "tcp://155.lchtime.com:1883" // 无ssl
	//mqtt_broker_address  = "ssl://155.lchtime.com:8883" // 支持ssl
	mqtt_broker_username = "1"
	mqtt_broker_password = "52399399"
)

var Client mqtt.Client
var heartOnce sync.Once

func MqInit(productKey, mac string) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqtt_broker_address)
	opts.SetClientID(misc.Mac())
	opts.SetUsername(mqtt_broker_username)
	opts.SetPassword(mqtt_broker_password)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	//	tlscfg, err := NewTLSConfig()
	//	if err != nil {
	//		panic(err)
	//	}
	//	opts.SetTLSConfig(tlscfg)

	opts.SetOnConnectHandler(func(cli mqtt.Client) {
		logs.Info("mqtt client connection success")
		chList := elink.ChannelSelectorList()
		for _, ch := range chList {
			s := fmt.Sprintf("%s/%s/%s/+/+/+/#", ch, productKey, mac)
			cli.Subscribe(s, 2, elinkmq.Handle)
		}
		heartOnce.Do(func() { time.AfterFunc(time.Second, HeartBeatStatus) })
	})

	opts.SetConnectionLostHandler(func(cli mqtt.Client, err error) {
		logs.Warn("mqtt client connection lost, ", err)
	})

	if out, err := jsoniter.Marshal(elinkmd.GatewayHeatbeats("", false)); err != nil {
		logs.Error("mqtt %s", err.Error())
	} else {
		opts.SetBinaryWill(
			fmt.Sprintf("data/0/%s/%s/patch/time", misc.Mac(), elinkmd.GatewayHeartbeat),
			out, 2, false)
	}
	Client = mqtt.NewClient(opts)

	elink.ManagaClient(true, elink.NewClient(elink.Hub, elinkmq.NewProvider(Client)))
	go Connect()
}
func NewTLSConfig() (*tls.Config, error) {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to
	// default openssl CA bundle.
	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM([]byte(cacert_pem))

	//	// Import client certificate/key pair
	//	cert, err := tls.X509KeyPair([]byte(cert_pem), []byte(key_pem))
	//	if err != nil {
	//		return nil, err
	//	}

	//	// Just to print out the client certificate..
	//	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	//	if err != nil {
	//		return nil, err
	//	}

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: true,
		// Certificates = list of certs client sends to server.
		//		Certificates: []tls.Certificate{cert},
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS12,
	}, nil
}

// 启动连接mqtt
func Connect() {
	logs.Info("mqtt client connecting...")
	if token := Client.Connect(); !token.WaitTimeout(time.Second*10) || (token.Error() != nil) {
		logs.Warn("mqtt client connect failed, ", token.Error())
		time.AfterFunc(time.Second*30, Connect)
	}
}

// 网关心跳包
func HeartBeatStatus() {
	defer time.AfterFunc(HeartBeatTime, HeartBeatStatus)
	if !Client.IsConnected() {
		return
	}

	// 心跳包推送
	func() {
		tp := elink.FormatPshSpecialTopic(ctrl.ChannelData,
			elinkmd.GatewayHeartbeat, elink.MethodPatch, elink.MessageTypeTime)
		out, err := jsoniter.Marshal(elinkmd.GatewayHeatbeats(tp, true))
		if err != nil {
			logs.Error("GatewayHeatbeats:", err)
			return
		}
		err = elink.Publish(tp, out)
		if err != nil {
			logs.Error("GatewayHeatbeats:", err)
		}
	}()

	// 系统监控信息推送
	func() {
		tp := elink.FormatPshSpecialTopic(ctrl.ChannelData,
			elinkmd.SystemMonitor, elink.MethodPatch, elink.MessageTypeTime)

		out, err := jsoniter.Marshal(elinkmd.GatewayMonitors(tp))
		if err != nil {
			logs.Error("GatewayMonitors:", err)
			return
		}
		err = elink.Publish(tp, out)
		if err != nil {
			logs.Error("GatewayHeatbeats:", err)
		}
	}()
}
