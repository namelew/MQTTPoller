package mqtt

import (
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/namelew/mqtt-bm-latency/logs"
)

type Client struct {
	Broker string
	ID     string
	KA     time.Duration
	Log    *logs.Log
	client mqtt.Client
}

func (c *Client) Create() {
	opts := mqtt.NewClientOptions().
		AddBroker(c.Broker).
		SetClientID(c.ID).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetKeepAlive(c.KA).
		SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {}).
		SetConnectionLostHandler(func(client mqtt.Client, reason error) {})

	client := mqtt.NewClient(opts)

	c.Log.Register("connect paho mqtt client to broker " + c.Broker)

	token := client.Connect()
	token.Wait()

	c.client = client
}

func (c *Client) Register(from string, qos int, quiet bool, handler mqtt.MessageHandler) {
	token := c.client.Subscribe(from, byte(qos), handler)
	token.Wait()
	if !quiet {
		c.Log.Register("Listen on topic " + from + " with qos " + strconv.Itoa(qos))
	}
}

func (c *Client) Unregister(from ...string) {
	for i := range from {
		token := c.client.Unsubscribe(from[i])
		token.Wait()
	}
}

func (c *Client) Send(to string, payload interface{}) {
	token := c.client.Publish(to, byte(1), false, payload)
	token.Wait()
}

func (c *Client) Disconnect(timeout uint) {
	c.client.Disconnect(timeout)
}
