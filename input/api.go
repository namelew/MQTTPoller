package input

type Info struct {
	Id            []int
	MemoryDisplay bool
	CpuDisplay    bool
	DiscDisplay   bool
}

type Experiment struct {
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

type Start struct {
	Id      []int
	Description Experiment
	ExeMode string
}