package messages

import (
	"encoding/json"
)

type CommandExperiment struct {
	Expid                 int64
	Attempts              int
	Tool                  string
	Broker                string
	Port                  int
	MqttVersion           int
	NumPublishers         int
	NumSubscriber         int
	QosPublisher          int
	QosSubscriber         int
	SharedSubscrition     bool
	Retain                bool
	Topic                 string
	Payload               int
	NumMessages           int
	RampUp                int
	RampDown              int
	Interval              int
	SubscriberTimeout     int
	Exec_time             int
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

type Info struct {
	Cpu  string
	Ram  uint64
	Disk uint64
}

type Command struct {
	Name        string                 `json:"name"`
	CommandType string                 `json:"commandType"`
	Arguments   map[string]interface{} `json:"arguments"`
}

type Worker struct {
	Id                  string
	Status              bool
	ReceiveConfirmation bool
	TestPing            bool
	Historic            ExperimentHistory
}

func (cmd *Command) ToCommandExperiment() *CommandExperiment {
	var cExp CommandExperiment

	data, err := json.Marshal(cmd.Arguments)

	if err != nil {
		return nil
	}

	err = json.Unmarshal(data, &cExp)

	if err != nil {
		return nil
	}

	return &cExp
}

func (cmdExp *CommandExperiment) Attach(cmd *Command) {
	data, err := json.Marshal(cmdExp)

	if err != nil {
		return
	}

	json.Unmarshal(data, &cmd.Arguments)
}
