package output

import "github.com/namelew/mqtt-bm-latency/messages"

type Worker struct {
	Id      int
	NetId   string
	Online  bool
	History []interface{}
}

type Info struct {
	Id    int
	NetId string
	Infos messages.Info
}