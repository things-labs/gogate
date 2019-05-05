package broad

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"time"

	"github.com/thinkgos/gogate/misc"
	"github.com/thinkgos/gogate/protocol/elinkmd"
	"github.com/thinkgos/gomo/elink"

	"github.com/astaxie/beego/logs"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	jsoniter "github.com/json-iterator/go"
)

const (
	mqttBrokerAddress = "tcp://155.lchtime.com:1883" // 无ssl
	//mqtt_broker_address  = "ssl://155.lchtime.com:8883" // 支持ssl
	mqttBrokerUsername = "1"
	mqttBrokerPassword = "52399399"
)

func NewMqClient(productKey, mac string) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttBrokerAddress)
	opts.SetClientID(misc.Mac())
	opts.SetUsername(mqttBrokerUsername)
	opts.SetPassword(mqttBrokerPassword)
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
			cli.Subscribe(s, 2, MessageHandle)
		}
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
	c := mqtt.NewClient(opts)
	connect := func() error {
		logs.Info("mqtt client connecting...")
		if token := c.Connect(); !token.WaitTimeout(time.Second*10) ||
			(token.Error() != nil) {
			logs.Warn("mqtt client connect failed, ", token.Error())
			return errors.New("mqtt client connect failed")
		}
		return nil
	}

	go func() {
		if connect() == nil {
			return
		}
		t := time.NewTimer(time.Second * 30)
		for {
			<-t.C
			if err := connect(); err != nil {
				t.Reset(time.Second * 30)
				continue
			}
			t.Stop()
			return
		}
	}()

	return c
}

func NewTLSConfig() (*tls.Config, error) {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to
	// default openssl CA bundle.
	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM([]byte(cacertPem))

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
