package messages

type CommandExperiment struct {
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

type SubscriberExperimentResult struct {
	Throughput        float64 `json:"max_throughput"`
	AvgThroughput     float64 `json:"avg_throughput"`
	ReceivedMessages  int `json:"received_messages"`
	PerSecondThrouput []int `json:"per_second_throungput"`
	Latency           float64 `json:"latency"`
	AvgLatency        float64 `json:"avg_latency"`
}

type PublisherExperimentResult struct {
	Throughput        float64 `json:"max_throughput"`
	AvgThroughput     float64 `json:"avg_throughput"`
	PubMessages       int `json:"publiqued_messages"`
	PerSecondThrouput []int `json:"per_second_throungput"`
}

type MetaExperimentResult struct {
	ID			uint64		`json:"id"`
	ExperimentError string `json:"error"`
	ToolName        string `json:"tool"`
	Literal         string `json:"literal"`
	LogFile			File `json:"log_file"`
}

type ExperimentResult struct {
	Meta      MetaExperimentResult `json:"meta"`
	Publish   PublisherExperimentResult `json:"publish"`
	Subscribe SubscriberExperimentResult `json:"subscribe"`
}

// type Tool struct {
// 	Name    string
// 	Version string
// }

type File struct{
	Name string `json:"name"`
	Data []byte `json:"data"`
	Extension string `json:"extension"`
}