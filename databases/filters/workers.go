package filters

type Worker struct {
	WorkerID uint64
	Token    string
	Online   bool
	Error    string
}
