package messages

import (
	"encoding/json"
)

type CommandExperiment struct {
	Expid                int64 		`json:"expid"`
	Attempts             int    	`json:"attempts"`
	Tool                 string 	`json:"tool"`
	Broker               string 	`json:"broker"`
	Port                 int		`json:"broker_port"`
	MqttVersion          int 		`json:"mqtt_version"`
	NumPublishers        int 		`json:"num_publishers"`
	NumSubscriber        int 		`json:"num_subscribers"`
	QosPublisher         int 		`json:"qos_publisher"`
	QosSubscriber        int 		`json:"qos_subscriber"`
	SharedSubscrition    bool 		`json:"shared_subscription"`
	Retain               bool 		`json:"retain"`
	Topic                string		`json:"topic"`
	Payload              int 		`json:"payload"`
	NumMessages          int 		`json:"num_messages"`
	RampUp               int 		`json:"ramp_up"`
	RampDown             int 		`json:"ramp_down"`
	Interval             int 		`json:"interval"`
	SubscriberTimeout    int 		`json:"subscriber_timeout"`
	ExecTime             int 		`json:"exec_time"`
	LogLevel             string 	`json:"log_level"`
	Ntp                  string 	`json:"ntp"`
	Output               bool 		`json:"output"`
	User                 string 	`json:"user_name"`
	Password             string 	`json:"password"`
	TlsTrustsore         string 	`json:"tls_truststore"`
	TlsTruststorePassword string 	`json:"tls_truststore_pass"`
	TlsKeystore          string 	`json:"tls_keystore"`
	TlsKeystorePassword  string 	`json:"tls_keystore_pass"`
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
