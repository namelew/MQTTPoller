package network

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/namelew/mqtt-bm-latency/packages/logs"
)

const (
	ID string = "Orquestrator"
	KEEPALIVE time.Duration = time.Second * 1000
)

var main_log *logs.Log = nil

func Create(b string, l *logs.Log) mqtt.Client{
	main_log = l

	opts := mqtt.NewClientOptions().
		AddBroker(b).
		SetClientID(ID).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetKeepAlive(KEEPALIVE).
		SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {}).
		SetConnectionLostHandler(func(client mqtt.Client, reason error) {})

	client := mqtt.NewClient(opts)

	main_log.Register("connect paho mqtt client to broker " + b)

	token := client.Connect()
	token.Wait()

	return client
}