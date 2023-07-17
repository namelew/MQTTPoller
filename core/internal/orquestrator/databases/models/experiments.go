package models

import "time"

type Experiment struct {
	ID                      uint64 `gorm:"primarykey"`
	ExperimentDeclarationID uint64
	Finish                  bool
	Error                   string
	CreatedAt               time.Time
	UpdatedAt               time.Time
	Workers                 []*Worker             `gorm:"many2many:experiments_workers;"`
	ExperimentDeclaration   ExperimentDeclaration `gorm:"foreignKey:ExperimentDeclarationID;references:ID"`
}

type ExperimentDeclaration struct {
	ID                    uint64 `gorm:"primarykey"`
	Tool                  string `json:"tool"`
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
}
