package models

import "time"

type Experiment struct {
	ID                      uint64 `gorm:"primarykey"`
	ExperimentDeclarationID uint64
	Finish                  bool
	Error                   string
	CreatedAt               time.Time
	UpdatedAt               time.Time
	Workers                 []*Worker `gorm:"many2many:experiments_workers;"`
}

type ExperimentDeclaration struct {
	ID                    uint64 `gorm:"primarykey"`
	Attempts              int
	Tool                  string
	Broker                string
	Port                  uint
	MqttVersion           uint8
	NumPublishers         uint
	NumSubscriber         uint
	QosPublisher          uint8
	QosSubscriber         uint8
	SharedSubscrition     bool
	Retain                bool
	Topic                 string
	Payload               uint64
	NumMessages           uint
	RampUp                int
	RampDown              int
	Interval              uint64
	SubscriberTimeout     uint64
	ExecTime              uint64
	LogLevel              string
	Ntp                   string
	Output                bool
	User                  string
	Password              string
	TlsTrustsore          string
	TlsTruststorePassword string
	TlsKeystore           string
	TlsKeystorePassword   string
}

type ExperimentResult struct {
	ID                  uint64 `gorm:"primarykey"`
	ExperimentID        uint64
	Error               string
	Tool                string
	Literal             string
	Filename            string
	FileData            []byte
	FileExtension       string
	SubThroughput       float64
	SubAvgThroughput    float64
	SubReceivedMessages uint
	SubLatency          float64
	SubAvgLatency       float64
	PubThroughput       float64
	PubAvgThroughput    float64
	PubMessages         uint
	Experiment          Experiment
}

type ExperimentResultPerSecondThrouput struct {
	ExperimentResultID uint64
	Time               uint64
	Value              int
	Action             string
}
