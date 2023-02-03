package output

type SubscriberExperimentResult struct {
	Throughput        float64 `json:"max_throughput"`
	AvgThroughput     float64 `json:"avg_throughput"`
	ReceivedMessages  int     `json:"received_messages"`
	PerSecondThrouput []int   `json:"per_second_throungput"`
	Latency           float64 `json:"latency"`
	AvgLatency        float64 `json:"avg_latency"`
}

type PublisherExperimentResult struct {
	Throughput        float64 `json:"max_throughput"`
	AvgThroughput     float64 `json:"avg_throughput"`
	PubMessages       int     `json:"publiqued_messages"`
	PerSecondThrouput []int   `json:"per_second_throungput"`
}

type MetaExperimentResult struct {
	ID              uint64 `json:"id"`
	ExperimentError string `json:"error"`
	ToolName        string `json:"tool"`
	Literal         string `json:"literal"`
	LogFile         File   `json:"log_file"`
}

type ExperimentResult struct {
	Meta      MetaExperimentResult       `json:"meta"`
	Publish   PublisherExperimentResult  `json:"publish"`
	Subscribe SubscriberExperimentResult `json:"subscribe"`
}

type File struct {
	Name      string `json:"name"`
	Data      []byte `json:"data"`
	Extension string `json:"extension"`
}