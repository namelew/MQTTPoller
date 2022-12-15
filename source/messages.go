package main

type commandExperiment struct {
	Expid                int64 `json:"expid"`
	Tool                 string `json:"tool"`
	Broker               string `json:"broker"`
	Port                 int `json:"broker_port"`
	MqttVersion          int `json:"mqtt_version"`
	NumPublishers        int `json:"num_publishers"`
	NumSubscriber        int `json:"num_subscribers"`
	QosPublisher         int `json:"qos_publisher"`
	QosSubscriber        int `json:"qos_subscriber"`
	SharedSubscrition    bool `json:"shared_subscription"`
	Retain               bool `json:"retain"`
	Topic                string `json:"topic"`
	Payload              int `json:"payload"`
	NumMessages          int `json:"num_messages"`
	RampUp               int `json:"ramp_up"`
	RampDown             int `json:"ramp_down"`
	Interval             int `json:"interval"`
	SubscriberTimeout    int `json:"subscriber_timeout"`
	Exec_time            int `json:"exec_time"`
	LogLevel             string `json:"log_level"`
	Ntp                  string `json:"ntp"`
	Output               bool `json:"output"`
	User                 string `json:"user_name"`
	Password             string `json:"password"`
	TlsTrustsore         string `json:"tls_truststore"`
	TlsTruststorePassword string `json:"tls_truststore_pass"`
	TlsKeystore          string `json:"tls_keystore"`
	TlsKeystorePassword  string `json:"tls_keystore_pass"`
}

type subscriberExperimentResult struct {
	Throughput        float64 `json:"max_throughput"`
	AvgThroughput     float64 `json:"avg_throughput"`
	ReceivedMessages  int `json:"received_messages"`
	PerSecondThrouput []int `json:"per_second_throungput"`
	Latency           float64 `json:"latency"`
	AvgLatency        float64 `json:"avg_latency"`
}

type publisherExperimentResult struct {
	Throughput        float64 `json:"max_throughput"`
	AvgThroughput     float64 `json:"avg_throughput"`
	PubMessages       int `json:"publiqued_messages"`
	PerSecondThrouput []int `json:"per_second_throungput"`
}

type metaExperimentResult struct {
	ID			uint64		`json:"id"`
	ExperimentError string `json:"error"`
	ToolName        string `json:"tool"`
	Literal         string `json:"literal"`
	LogFile			file `json:"log_file"`
}

type experimentResult struct {
	Meta      metaExperimentResult `json:"meta"`
	Publish   publisherExperimentResult `json:"publish"`
	Subscribe subscriberExperimentResult `json:"subscribe"`
}

// type tool struct {
// 	Name    string
// 	Version string
// }

type info struct {
	MemoryDisplay bool
	CpuDisplay    bool
	DiscDisplay   bool
}

type infoDisplay struct{
	Cpu 	string
	Ram		uint64
	Disk 	uint64 
}

type command struct{
	Name string `json:"name"`
	Type string `json:"type"`
	Arguments map[string]interface{} `json:"arguments"`
}

type file struct{
	Name string `json:"name"`
	Data []byte `json:"data"`
	Extension string `json:"extension"`
}

type status struct{
	Type string `json:"type"`
	Status string `json:"status"`
	Attr command `json:"attr"`
}

type session struct{
	Id int
	Finish bool
	Status status
}