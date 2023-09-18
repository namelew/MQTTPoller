package filters

type Worker struct {
	WorkerID uint64 `json:"id"`
	Token    string `json:"token"`
	Online   bool   `json:"online"`
	Error    string `json:"error"`
}
