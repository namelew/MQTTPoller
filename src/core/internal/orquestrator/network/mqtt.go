package network

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/namelew/mqtt-poller/src/core/packages/logs"
)

const (
	ID        string        = "Orquestrator"
	KEEPALIVE time.Duration = time.Second * 30
)

var main_log *logs.Log = nil

func Create(b string, l *logs.Log) mqtt.Client {
	main_log = l

	opts := mqtt.NewClientOptions().
		AddBroker(b).
		SetClientID(ID).
		SetCleanSession(false).
		SetAutoReconnect(true).
		SetKeepAlive(KEEPALIVE).
		SetConnectionLostHandler(func(client mqtt.Client, reason error) {
			l.Register("Connection lost. Reason: " + reason.Error())
		})

	client := mqtt.NewClient(opts)

	main_log.Register("connect paho mqtt client to broker " + b)

	token := client.Connect()
	token.Wait()

	return client
}
