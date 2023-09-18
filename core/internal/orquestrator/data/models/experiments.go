package models

import "github.com/namelew/mqtt-poller/core/packages/messages"

type Experiment struct {
	ID                    uint64
	Broker                string `json:"broker"`
	Port                  int    `json:"port"`
	MqttVersion           int    `json:"mqttVersion"`
	NumPublishers         int    `json:"numPublishers"`
	NumSubscriber         int    `json:"numSubscribers"`
	QosPublisher          int    `json:"qosPublisher"`
	QosSubscriber         int    `json:"qosSubscriber"`
	SharedSubscrition     bool   `json:"sharedSubscription"`
	Retain                bool   `json:"retain"`
	Topic                 string `json:"topic"`
	Payload               int    `json:"payload"`
	NumMessages           int    `json:"numMessages"`
	RampUp                int    `json:"ramUp"`
	RampDown              int    `json:"rampDown"`
	Interval              int    `json:"interval"`
	SubscriberTimeout     int    `json:"subscriberTimeout"`
	ExecTime              int    `json:"execTime"`
	LogLevel              string `json:"logLevel"`
	Ntp                   string `json:"ntp"`
	Output                bool   `json:"output"`
	User                  string `json:"username"`
	Password              string `json:"password"`
	TlsTrustsore          string `json:"tlsTruststore"`
	TlsTruststorePassword string `json:"tlsTruststorePass"`
	TlsKeystore           string `json:"tlsKeystore"`
	TlsKeystorePassword   string `json:"tlsKeystorePass"`
	Finish                bool
	Error                 string
	Results               []messages.ExperimentResult
	WorkerIDs             []string
}
